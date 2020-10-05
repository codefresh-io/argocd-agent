package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/installer/pkg/cliconfig"
	"github.com/codefresh-io/argocd-listener/installer/pkg/fs"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os/user"
	"path"
	"regexp"
	"strconv"
)

var installCmdOptions struct {
	kube struct {
		namespace    string
		inCluster    bool
		context      string
		nodeSelector string
		configPath   string
	}
	Argo struct {
		Host     string
		Username string
		Password string
		Token    string
	}
	Codefresh struct {
		Host        string
		Token       string
		Integration string
		AutoSync    string
	}
	Agent struct {
		Version		string
	}
}

func ensureIntegration() error {
	err := holder.ApiHolder.CreateIntegration(installCmdOptions.Codefresh.Integration, installCmdOptions.Argo.Host, installCmdOptions.Argo.Username, installCmdOptions.Argo.Password, installCmdOptions.Argo.Token)
	if err == nil {
		return nil
	}

	codefreshErr, ok := err.(*codefresh.CodefreshError)
	if !ok {
		return err
	}

	if codefreshErr.Status != 409 {
		return codefreshErr
	}

	err, needUpdate := prompt.Confirm("You already have integration with this name, do you want to update it")
	if err != nil {
		return err
	}

	if !needUpdate {
		return fmt.Errorf("you should update integration")
	}

	err = holder.ApiHolder.UpdateIntegration(installCmdOptions.Codefresh.Integration, installCmdOptions.Argo.Host, installCmdOptions.Argo.Username, installCmdOptions.Argo.Password, installCmdOptions.Argo.Token)

	if err != nil {
		return err
	}

	return nil
}

func sendArgoAgentInstalledEvent(status string, reason string) {
	props := make(map[string]string)
	props["status"] = status
	props["reason"] = reason
	_ = holder.ApiHolder.SendEvent("agent.installed", props)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install agent",
	Long:  `Install agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if installCmdOptions.Codefresh.Token == "" || installCmdOptions.Codefresh.Host == "" {
			config, err := cliconfig.GetCurrentConfig()
			if err != nil {
				return err
			}
			installCmdOptions.Codefresh.Token = config.Token
			installCmdOptions.Codefresh.Host = config.Url
		}

		err = prompt.InputWithDefault(&installCmdOptions.Codefresh.Integration, "Codefresh integration name", "argocd")
		if err != nil {
			return err
		}

		err = prompt.Input(&installCmdOptions.Argo.Host, "Argo host, example: https://example.com")
		if err != nil {
			return err
		}

		withProtocol, err := regexp.MatchString("^https?://", installCmdOptions.Argo.Host)
		if err != nil {
			return err
		}

		// customer not put protocol during installation
		if !withProtocol {
			installCmdOptions.Argo.Host = "https://" + installCmdOptions.Argo.Host
		}

		// removing / in the end
		installCmdOptions.Argo.Host = regexp.MustCompile("/+$").ReplaceAllString(installCmdOptions.Argo.Host, "")

		err, useArgocdToken := prompt.Confirm("Do you want use argocd auth token instead username/password auth?")
		if err != nil {
			return err
		}

		if useArgocdToken {
			err = prompt.InputWithDefault(&installCmdOptions.Argo.Token, "Argo token", "")
			if err != nil {
				return err
			}
		} else {
			err = prompt.InputWithDefault(&installCmdOptions.Argo.Username, "Argo username", "admin")
			if err != nil {
				return err
			}

			err = prompt.InputPassword(&installCmdOptions.Argo.Password, "Argo password")
			if err != nil {
				return err
			}
		}

		holder.ApiHolder = codefresh.Api{
			Token:       installCmdOptions.Codefresh.Token,
			Host:        installCmdOptions.Codefresh.Host,
			Integration: installCmdOptions.Codefresh.Integration,
		}

		err = ensureIntegration()
		if err != nil {
			sendArgoAgentInstalledEvent(FAILED, err.Error())
			return err
		}

		kubeConfigPath := installCmdOptions.kube.configPath
		kubeOptions := installCmdOptions.kube

		if kubeOptions.context == "" {
			contexts, err := kube.GetAllContexts(kubeConfigPath)
			if err != nil {
				return err
			}

			err, selectedContext := prompt.Select(contexts, "Select Kubernetes context")
			kubeOptions.context = selectedContext
		}

		kubeClient, err := kube.New(&kube.Options{
			ContextName:      kubeOptions.context,
			Namespace:        kubeOptions.namespace,
			PathToKubeConfig: kubeConfigPath,
			InCluster:        kubeOptions.inCluster,
		})
		if err != nil {
			return err
		}

		namespaces, err := kubeClient.GetNamespaces()
		if err != nil {
			err = prompt.InputWithDefault(&kubeOptions.namespace, "Kubernetes namespace to install", "default")
			if err != nil {
				return err
			}
		} else {
			err, selectedNamespace := prompt.Select(namespaces, "Select Kubernetes namespace")
			if err != nil {
				return err
			}
			kubeOptions.namespace = selectedNamespace
		}

		err, autoSync := prompt.Confirm("Do you want auto sync argo apps to codefresh?")
		if err != nil {
			return err
		}

		installCmdOptions.Codefresh.AutoSync = strconv.FormatBool(autoSync)

		installCmdOptions.Codefresh.Token = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Codefresh.Token))
		installCmdOptions.Argo.Token = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Argo.Token))

		installOptions := templates.InstallOptions{
			Templates:      kubernetes.TemplatesMap(),
			TemplateValues: structs.Map(installCmdOptions),
			Namespace:      kubeOptions.namespace,
			KubeClientSet:  kubeClient.GetClientSet(),
		}

		var kind, name string
		err, kind, name = templates.Install(&installOptions)

		if err != nil {
			msg := fmt.Sprintf("Argo agent installation resource \"%s\" with name \"%s\" finished with error , reason: %v ", kind, name, err)
			sendArgoAgentInstalledEvent(FAILED, msg)
			return errors.New(msg)
		}

		sendArgoAgentInstalledEvent(SUCCESS, "")

		logger.Success(fmt.Sprintf("Argo agent installation finished successfully to namespace \"%s\"", kubeOptions.namespace))

		return nil
	},
}

func init() {

	rootCmd.AddCommand(installCmd)
	flags := installCmd.Flags()

	flags.StringVar(&installCmdOptions.Agent.Version, "agent-version", fs.GetAgentVersion(), "")

	flags.StringVar(&installCmdOptions.Argo.Host, "argo-host", "", "")
	flags.StringVar(&installCmdOptions.Argo.Username, "argo-username", "", "")
	flags.StringVar(&installCmdOptions.Argo.Password, "argo-password", "", "")

	flags.StringVar(&installCmdOptions.Codefresh.Host, "codefresh-host", "http://local.codefresh.io", "")
	flags.StringVar(&installCmdOptions.Codefresh.Token, "codefresh-token", "", "")
	flags.StringVar(&installCmdOptions.Codefresh.Integration, "codefresh-integration", "", "")

	flags.StringVar(&installCmdOptions.kube.namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which Argo agent should be installed [$KUBE_NAMESPACE]")
	flags.StringVar(&installCmdOptions.kube.context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which Argo agent should be installed (default is current-context) [$KUBE_CONTEXT]")
	flags.BoolVar(&installCmdOptions.kube.inCluster, "in-cluster", false, "Set flag if Argo agent is been installed from inside a cluster")

	var kubeConfigPath string
	currentUser, _ := user.Current()
	if currentUser != nil {
		kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
	}

	flags.StringVar(&installCmdOptions.kube.configPath, "kubeconfig", kubeConfigPath, "Path to kubeconfig for retrieve contexts")

}

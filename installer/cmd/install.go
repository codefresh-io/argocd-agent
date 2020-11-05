package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/questionnaire"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/codefresh-io/argocd-listener/installer/pkg/util"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	"os/user"
	"path"
)

var installCmdOptions = install.InstallCmdOptions{}

// variable derived from ldflag
var version = ""

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

	needUpdate := installCmdOptions.Argo.Update
	if !needUpdate {
		err, needUpdate = prompt.Confirm("You already have integration with this name, do you want to update it")
		if err != nil {
			return err
		}
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
		logger.Success("This installer will guide you through the Codefresh ArgoCD installation agent to integrate your ArgoCD with Codefresh")

		kubeConfigPath := installCmdOptions.Kube.ConfigPath
		kubeOptions := installCmdOptions.Kube

		_, cluster := questionnaire.AskAboutKubeContext(&installCmdOptions)

		kubeClient, err := kube.New(&kube.Options{
			ContextName:      kubeOptions.Context,
			Namespace:        kubeOptions.Namespace,
			PathToKubeConfig: kubeConfigPath,
			InCluster:        kubeOptions.InCluster,
		})
		if err != nil {
			return err
		}
		_ = questionnaire.AskAboutNamespace(&installCmdOptions, kubeClient)

		kubeOptions = installCmdOptions.Kube

		argoServerPlaced := kubeClient.IsArgoServerOnCluster(kubeOptions.Namespace)
		if !argoServerPlaced {
			msg := fmt.Sprintf("We didnt find ArgoCD on \"%s/%s\"", cluster, kubeOptions.Namespace)
			sendArgoAgentInstalledEvent(FAILED, msg)
			return errors.New(msg)
		}

		_ = questionnaire.AskAboutCodefreshCredentials(&installCmdOptions)

		err = prompt.InputWithDefault(&installCmdOptions.Codefresh.Integration, "Codefresh integration name", "argocd")
		if err != nil {
			return err
		}

		_ = questionnaire.AskAboutArgoCredentials(&installCmdOptions)

		holder.ApiHolder = codefresh.Api{
			Token:       installCmdOptions.Codefresh.Token,
			Host:        installCmdOptions.Codefresh.Host,
			Integration: installCmdOptions.Codefresh.Integration,
		}

		_ = questionnaire.AskAboutGitContext(&installCmdOptions)

		err = ensureIntegration()
		if err != nil {
			sendArgoAgentInstalledEvent(FAILED, err.Error())
			return err
		}

		// Need check if we want support not in cluster mode with Product owner
		installCmdOptions.Kube.InCluster = true

		questionnaire.AskAboutSyncOptions(&installCmdOptions)

		installCmdOptions.Codefresh.Token = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Codefresh.Token))
		installCmdOptions.Argo.Token = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Argo.Token))

		installOptions := templates.InstallOptions{
			Templates:        kubernetes.TemplatesMap(),
			TemplateValues:   structs.Map(installCmdOptions),
			Namespace:        kubeOptions.Namespace,
			KubeClientSet:    kubeClient.GetClientSet(),
			KubeManifestPath: installCmdOptions.Kube.ManifestPath,
		}

		var kind, name string
		err, kind, name = templates.Install(&installOptions)

		if err != nil {
			msg := fmt.Sprintf("Argo agent installation resource \"%s\" with name \"%s\" finished with error , reason: %v ", kind, name, err)
			sendArgoAgentInstalledEvent(FAILED, msg)
			return errors.New(msg)
		}

		sendArgoAgentInstalledEvent(SUCCESS, "")

		logger.Success(fmt.Sprintf("Argo agent installation finished successfully to namespace \"%s\"", kubeOptions.Namespace))

		return nil
	},
}

func init() {

	rootCmd.AddCommand(installCmd)
	flags := installCmd.Flags()

	flags.StringVar(&installCmdOptions.Agent.Version, "agent-version", util.ResolvePackageVersion(version), "")
	flags.StringVar(&installCmdOptions.Argo.Host, "argo-host", "", "")
	flags.StringVar(&installCmdOptions.Argo.Token, "argo-token", "", "")
	flags.StringVar(&installCmdOptions.Argo.Username, "argo-username", "", "")
	flags.StringVar(&installCmdOptions.Argo.Password, "argo-password", "", "")
	flags.BoolVar(&installCmdOptions.Argo.Update, "update", false, "Update integration if exists")

	flags.StringVar(&installCmdOptions.Codefresh.Host, "codefresh-host", "http://local.codefresh.io", "")
	flags.StringVar(&installCmdOptions.Codefresh.Token, "codefresh-token", "", "")
	flags.StringVar(&installCmdOptions.Codefresh.Integration, "codefresh-integration", "", "Argocd integration in Codefresh")
	flags.StringVar(&installCmdOptions.Codefresh.SyncMode, "sync-mode", "", "")
	flags.StringArrayVar(&installCmdOptions.Codefresh.ApplicationsForSyncArr, "sync-apps", make([]string, 0), "")

	flags.StringVar(&installCmdOptions.Kube.ManifestPath, "kube-manifest-path", "", "")
	flags.StringVar(&installCmdOptions.Kube.Namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which Argo agent should be installed [$KUBE_NAMESPACE]")
	flags.StringVar(&installCmdOptions.Kube.Context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which Argo agent should be installed (default is current-context) [$KUBE_CONTEXT]")
	flags.BoolVar(&installCmdOptions.Kube.InCluster, "in-cluster", false, "Set flag if Argo agent is been installed from inside a cluster")

	flags.StringVar(&installCmdOptions.Git.Integration, "git-integration", "", "Name of git integration in Codefresh")

	var kubeConfigPath string
	currentUser, _ := user.Current()
	if currentUser != nil {
		kubeConfigPath = os.Getenv("KUBECONFIG")
		if kubeConfigPath == "" {
			kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
		}
	}

	flags.StringVar(&installCmdOptions.Kube.ConfigPath, "kubeconfig", kubeConfigPath, "Path to kubeconfig for retrieve contexts")

}

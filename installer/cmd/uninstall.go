package cmd

import (
	"errors"
	"fmt"
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
)

var uninstallCmdOptions struct {
	kube struct {
		namespace  string
		inCluster  bool
		context    string
		configPath string
	}
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall agent",
	Long:  `Uninstall agent`,
	RunE: func(cmd *cobra.Command, args []string) error {

		kubeConfigPath := installCmdOptions.kube.configPath
		kubeOptions := uninstallCmdOptions.kube

		if uninstallCmdOptions.kube.context == "" {
			contexts, err := kube.GetAllContexts(kubeConfigPath)
			if err != nil {
				return err
			}

			err, selectedContext := prompt.Select(contexts, "Select Kubernetes context")
			if err != nil {
				return err
			}
			kubeOptions.context = selectedContext
		}

		err := prompt.InputWithDefault(&kubeOptions.namespace, "Kubernetes namespace to uninstall", "default")
		if err != nil {
			return err
		}

		kubeClient, err := kube.New(&kube.Options{
			ContextName:      kubeOptions.context,
			Namespace:        kubeOptions.namespace,
			PathToKubeConfig: kubeConfigPath,
			InCluster:        kubeOptions.inCluster,
		})

		if err != nil {
			panic(err)
		}

		uninstallOptions := templates.DeleteOptions{
			Templates:      kubernetes.TemplatesMap(),
			TemplateValues: structs.Map(uninstallCmdOptions),
			Namespace:      kubeOptions.namespace,
			KubeClientSet:  kubeClient.GetClientSet(),
		}

		var kind, name string
		err, kind, name = templates.Delete(&uninstallOptions)

		if err != nil {
			return errors.New(fmt.Sprintf("Argo agent uninstallation resource \"%s\" with name \"%s\" finished with error , reason: %v ", kind, name, err))
		}

		logger.Success(fmt.Sprintf("Argo agent uninstallation finished successfully to namespace \"%s\"", kubeOptions.namespace))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().StringVar(&uninstallCmdOptions.kube.namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which venona should be installed [$KUBE_NAMESPACE]")
	uninstallCmd.Flags().StringVar(&uninstallCmdOptions.kube.context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which venona should be installed (default is current-context) [$KUBE_CONTEXT]")
	uninstallCmd.Flags().BoolVar(&uninstallCmdOptions.kube.inCluster, "in-cluster", false, "Set flag if venona is been installed from inside a cluster")

	var kubeConfigPath string
	currentUser, _ := user.Current()
	if currentUser != nil {
		kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
	}

	uninstallCmd.Flags().StringVar(&installCmdOptions.kube.configPath, "kubeconfig", kubeConfigPath, "Path to kubeconfig for retrieve contexts")
}

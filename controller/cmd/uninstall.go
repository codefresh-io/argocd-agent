package cmd

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/controller/pkg/install"
	"github.com/codefresh-io/argocd-listener/controller/pkg/kube"
	"github.com/codefresh-io/argocd-listener/controller/pkg/logger"
	"github.com/codefresh-io/argocd-listener/controller/pkg/questionnaire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path"
)

var uninstallCmdOptions = install.CmdOptions{}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall gitops codefresh",
	Long:  `Uninstall gitops codefresh`,
	RunE: func(cmd *cobra.Command, args []string) error {

		_ = questionnaire.AskAboutKubeContext(&uninstallCmdOptions)
		kubeOptions := uninstallCmdOptions.Kube
		kubeClient, err := kube.New(&kube.Options{
			ContextName:      kubeOptions.Context,
			Namespace:        kubeOptions.Namespace,
			PathToKubeConfig: kubeOptions.ConfigPath,
		})

		if err != nil {
			return failUninstall(fmt.Sprintf("Can't create kube client: \"%s\"", err.Error()))
		}

		_ = questionnaire.AskAboutNamespace(&uninstallCmdOptions, kubeClient)
		_ = questionnaire.AskAboutManifest(&uninstallCmdOptions)
		err = kubeClient.DeleteObjects(uninstallCmdOptions.Kube.ManifestPath)
		if err != nil {
			return failUninstall(fmt.Sprintf("Can't delete kube objects: \"%s\"", err.Error()))
		}

		//sendArgoAgentUninstalledEvent(SUCCESS, "")

		logger.Success(fmt.Sprintf("Codefresh gitops controller uninstallation finished successfully"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	flags := uninstallCmd.Flags()

	flags.StringVar(&uninstallCmdOptions.Kube.Namespace, "kube-namespace", "argocd", "Namespace in Kubernetes cluster")
	flags.StringVar(&uninstallCmdOptions.Kube.ManifestPath, "install-manifest", "", "Url of argocd install manifest")

	var kubeConfigPath string
	currentUser, _ := user.Current()
	if currentUser != nil {
		kubeConfigPath = os.Getenv("KUBECONFIG")
		if kubeConfigPath == "" {
			kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
		}
	}

	flags.StringVar(&uninstallCmdOptions.Kube.Context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which Argo agent should be installed (default is current-context) [$KUBE_CONTEXT]")
	flags.StringVar(&uninstallCmdOptions.Kube.ConfigPath, "kube-config-path", kubeConfigPath, "Path to kubeconfig file (default is $HOME/.kube/config)")
}

func failUninstall(msg string) error {
	sendControllerInstalledEvent(FAILED, msg)
	return errors.New(msg)
}

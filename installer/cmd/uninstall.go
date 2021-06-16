package cmd

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/helper"
	"github.com/codefresh-io/argocd-listener/installer/pkg/uninstall"
	"github.com/codefresh-io/argocd-listener/installer/pkg/uninstall/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	"os/user"
)

var uninstallCmdOptions = uninstall.CmdOptions{}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall agent",
	Long:  `Uninstall agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		uninstallHandler := handler.New(uninstallCmdOptions)
		err := uninstallHandler.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().StringVar(&uninstallCmdOptions.Kube.Namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which venona should be installed [$KUBE_NAMESPACE]")
	uninstallCmd.Flags().StringVar(&uninstallCmdOptions.Kube.Context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which venona should be installed (default is current-context) [$KUBE_CONTEXT]")
	uninstallCmd.Flags().BoolVar(&uninstallCmdOptions.Kube.InCluster, "in-cluster", false, "Set flag if Argo agent is been installed from inside a cluster")

	var kubeConfigPath string
	currentUser, _ := user.Current()
	if currentUser != nil {
		kubeConfigPath = os.Getenv("KUBECONFIG")
		if kubeConfigPath == "" {
			kubeConfigPath = helper.GetDefaultKubeConfigPath(currentUser.HomeDir)
		}
	}

	uninstallCmd.Flags().StringVar(&uninstallCmdOptions.Kube.ConfigPath, "kubeconfig", kubeConfigPath, "Path to kubeconfig for retrieve contexts")
}

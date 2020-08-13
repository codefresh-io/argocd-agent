package cmd

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
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
		namespace string
		inCluster bool
		context   string
	}
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall agent",
	Long:  `Uninstall agent`,
	Run: func(cmd *cobra.Command, args []string) {

		var kubeConfigPath string
		currentUser, _ := user.Current()
		if currentUser != nil {
			kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
		}

		kubeOptions := uninstallCmdOptions.kube

		cs, err := kube.ClientBuilder(kubeOptions.context, kubeOptions.namespace, kubeConfigPath, kubeOptions.inCluster).BuildClient()

		if err != nil {
			panic(err)
		}

		uninstallOptions := templates.DeleteOptions{
			Templates:      kubernetes.TemplatesMap(),
			TemplateValues: structs.Map(uninstallCmdOptions),
			Namespace:      kubeOptions.namespace,
			KubeClientSet:  cs,
		}

		templates.Delete(&uninstallOptions)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().StringVar(&uninstallCmdOptions.kube.namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which venona should be installed [$KUBE_NAMESPACE]")
	uninstallCmd.Flags().StringVar(&uninstallCmdOptions.kube.context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which venona should be installed (default is current-context) [$KUBE_CONTEXT]")
	uninstallCmd.Flags().BoolVar(&uninstallCmdOptions.kube.inCluster, "in-cluster", false, "Set flag if venona is been installed from inside a cluster")

}

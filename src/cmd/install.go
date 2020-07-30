package cmd

import (
	"github.com/codefresh-io/argocd-listener/src/pkg/kube"
	"github.com/codefresh-io/argocd-listener/src/pkg/templates"
	"github.com/codefresh-io/argocd-listener/src/pkg/templates/kubernetes"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
)

var installCmdOptions struct {
	ClusterName            string
	ClusterNamespace       string
	clusterNameInCodefresh string
	kube                   struct {
		namespace    string
		inCluster    bool
		context      string
		nodeSelector string
	}
	Argo struct {
		Host     string
		Username string
		Password string
	}
	Codefresh struct {
		Host  string
		Token string
	}
	Namespace string
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {

		installCmdOptions.Namespace = "default"

		cs, _ := kube.ClientBuilder("argocd", "default", "/Users/pashavictorovich/.kube/config", false).BuildClient()
		installOptions := templates.InstallOptions{
			Templates:      kubernetes.TemplatesMap(),
			TemplateValues: structs.Map(installCmdOptions),
			Namespace:      "default",
			KubeClientSet:  cs,
		}

		templates.Install(&installOptions)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&installCmdOptions.Argo.Host, "argo-host", "https://34.71.103.174/", "")
	installCmd.Flags().StringVar(&installCmdOptions.Argo.Username, "argo-username", "admin", "")
	installCmd.Flags().StringVar(&installCmdOptions.Argo.Password, "argo-password", "newpassword", "")

	installCmd.Flags().StringVar(&installCmdOptions.Codefresh.Host, "codefresh-host", "https://g.codefresh.io", "")
	installCmd.Flags().StringVar(&installCmdOptions.Codefresh.Token, "codefresh-token", "test", "")

	//installCmd.Flags().StringVar(&installCmdOptions.kube.namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which venona should be installed [$KUBE_NAMESPACE]")
	//installCmd.Flags().StringVar(&installCmdOptions.kube.context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which venona should be installed (default is current-context) [$KUBE_CONTEXT]")

}

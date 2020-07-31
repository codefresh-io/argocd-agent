package cmd

import (
	"github.com/codefresh-io/argocd-listener/src/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/src/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/src/installer/pkg/templates/kubernetes"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os/user"
	"path"
)

var installCmdOptions struct {
	kube struct {
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

		var kubeConfigPath string
		currentUser, _ := user.Current()
		if currentUser != nil {
			kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
		}

		kubeOptions := installCmdOptions.kube

		cs, err := kube.ClientBuilder(kubeOptions.context, kubeOptions.namespace, kubeConfigPath, kubeOptions.inCluster).BuildClient()

		if err != nil {
			panic(err)
		}

		installOptions := templates.InstallOptions{
			Templates:      kubernetes.TemplatesMap(),
			TemplateValues: structs.Map(installCmdOptions),
			Namespace:      kubeOptions.namespace,
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

	installCmd.Flags().StringVar(&installCmdOptions.kube.namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which venona should be installed [$KUBE_NAMESPACE]")
	installCmd.Flags().StringVar(&installCmdOptions.kube.context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which venona should be installed (default is current-context) [$KUBE_CONTEXT]")
	installCmd.Flags().BoolVar(&installCmdOptions.kube.inCluster, "in-cluster", false, "Set flag if venona is been installed from inside a cluster")

}

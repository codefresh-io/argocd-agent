package cmd

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/update"
	"github.com/codefresh-io/argocd-listener/installer/pkg/update/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	"os/user"
	"path"
)

var updateCmdOptions = update.CmdOptions{}

var updateCMD = &cobra.Command{
	Use:   "update",
	Short: "Update agent",
	Long:  `Update agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		updateHandler := handler.New(updateCmdOptions, version)
		err := updateHandler.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(updateCMD)
	flags := updateCMD.Flags()

	flags.StringVar(&updateCmdOptions.Kube.Namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which Argo agent should be updated [$KUBE_NAMESPACE]")
	flags.StringVar(&updateCmdOptions.Kube.Context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which Argo agent should be updated (default is current-context) [$KUBE_CONTEXT]")
	flags.BoolVar(&updateCmdOptions.Kube.InCluster, "in-cluster", false, "Set flag if Argo agent is been updated from inside a cluster")
	flags.StringVar(&updateCmdOptions.Codefresh.Suffix, "codefresh-agent-suffix", "", "Agent suffix")

	var kubeConfigPath string
	currentUser, _ := user.Current()
	if currentUser != nil {
		kubeConfigPath = os.Getenv("KUBECONFIG")
		if kubeConfigPath == "" {
			kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
		}
	}

	flags.StringVar(&updateCmdOptions.Kube.ConfigPath, "kubeconfig", kubeConfigPath, "Path to kubeconfig for retrieve contexts")
}

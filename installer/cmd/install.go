package cmd

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/fs"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/helper"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	"os/user"
)

var installCmdOptions = entity.InstallCmdOptions{}

// variable derived from ldflag
var version = ""

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install agent",
	Long:  `Install agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Success("This installer will guide you through the Codefresh ArgoCD installation agent to integrate your ArgoCD with Codefresh")
		if installCmdOptions.Codefresh.Suffix != "" {
			installCmdOptions.Codefresh.Suffix = "-" + installCmdOptions.Codefresh.Suffix
		}

		err, manifest := install.Run(installCmdOptions)
		if err != nil {
			return err
		}

		if installCmdOptions.Kube.ManifestPath != "" {
			err = fs.WriteFile(installCmdOptions.Kube.ManifestPath, manifest)
		}
		return err
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
	flags.StringVar(&installCmdOptions.Codefresh.Suffix, "codefresh-agent-suffix", "", "Agent suffix")
	flags.StringVar(&installCmdOptions.Codefresh.SyncMode, "sync-mode", "", "")
	flags.StringArrayVar(&installCmdOptions.Codefresh.ApplicationsForSyncArr, "sync-apps", make([]string, 0), "")
	flags.StringVar(&installCmdOptions.Codefresh.Provider, "Provider", "argocd", "")

	flags.StringVar(&installCmdOptions.Kube.ManifestPath, "output", "", "Path to k8s manifest output file, example: /home/user/out.yaml")
	flags.StringVar(&installCmdOptions.Kube.Namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which Argo agent should be installed [$KUBE_NAMESPACE]")
	flags.StringVar(&installCmdOptions.Kube.Context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which Argo agent should be installed (default is current-context) [$KUBE_CONTEXT]")
	flags.BoolVar(&installCmdOptions.Kube.InCluster, "in-cluster", false, "Set flag if ArgoCD agent is been installed from inside a cluster")

	flags.StringVar(&installCmdOptions.Git.Integration, "git-integration", "", "Name of git integration in Codefresh")

	flags.StringVar(&installCmdOptions.Host.HttpProxy, "http-proxy", "", "Http proxy")
	flags.StringVar(&installCmdOptions.Host.HttpsProxy, "https-proxy", "", "Https proxy")
	flags.StringVar(&installCmdOptions.NewRelic.Key, "new-relic", "", "")
	flags.StringVar(&installCmdOptions.Env.Name, "env-name", "kubernetes", "")

	flags.IntVar(&installCmdOptions.Replicas, "replicas", 1, "Amount of replicas for stateful set")

	var kubeConfigPath string
	currentUser, _ := user.Current()
	if currentUser != nil {
		kubeConfigPath = os.Getenv("KUBECONFIG")
		if kubeConfigPath == "" {
			kubeConfigPath = helper.GetDefaultKubeConfigPath(currentUser.HomeDir)
		}
	}

	flags.StringVar(&installCmdOptions.Kube.ConfigPath, "kubeconfig", kubeConfigPath, "Path to kubeconfig for retrieve contexts")

}

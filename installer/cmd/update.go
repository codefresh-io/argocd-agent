package cmd

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/obj/kubeobj"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os/user"
	"path"
)

var updateCmdOptions struct {
	kube struct {
		namespace  string
		inCluster  bool
		context    string
		configPath string
	}
}

func updateDeploymentWithNewVersion(clientSet *kubernetes.Clientset, namespace string) error {
	deploymentList, err := kubeobj.GetDeployments(clientSet, namespace, "app=cf-argocd-agent")

	if err != nil {
		return errors.New(fmt.Sprintf("Argo agent update finished with error , reason: %v ", err))
	}

	if len(deploymentList.Items) == 0 {
		return errors.New("Argo agent failed to update because no deployments were found")
	}

	deployment := &deploymentList.Items[0]

	envs := deployment.Spec.Template.Spec.Containers[0].Env

	newEnvs := make([]v1.EnvVar, 0)

	for _, env := range envs {
		if env.Name == "AGENT_VERSION" {
			env.Value = util.ResolvePackageVersion("")
		}
		newEnvs = append(newEnvs, env)
	}

	deployment.Spec.Template.Spec.Containers[0].Env = newEnvs

	_, err = kubeobj.UpdateDeployment(clientSet, deployment, namespace)

	return err
}

var updateCMD = &cobra.Command{
	Use:   "update",
	Short: "Update agent",
	Long:  `Update agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		kubeConfigPath := installCmdOptions.Kube.ConfigPath
		kubeOptions := updateCmdOptions.kube

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
			err = prompt.InputWithDefault(&kubeOptions.namespace, "Kubernetes namespace to update", "default")
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

		err = updateDeploymentWithNewVersion(kubeClient.GetClientSet(), kubeOptions.namespace)

		if err != nil {
			return errors.New(fmt.Sprintf("Argo agent update finished with error , reason: %v ", err))
		}

		logger.Success(fmt.Sprintf("Argo agent update finished successfully to namespace \"%s\"", kubeOptions.namespace))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCMD)
	flags := updateCMD.Flags()

	flags.StringVar(&updateCmdOptions.kube.namespace, "kube-namespace", viper.GetString("kube-namespace"), "Name of the namespace on which Argo agent should be updated [$KUBE_NAMESPACE]")
	flags.StringVar(&updateCmdOptions.kube.context, "kube-context-name", viper.GetString("kube-context"), "Name of the kubernetes context on which Argo agent should be updated (default is current-context) [$KUBE_CONTEXT]")
	flags.BoolVar(&updateCmdOptions.kube.inCluster, "in-cluster", false, "Set flag if Argo agent is been updated from inside a cluster")

	var kubeConfigPath string
	currentUser, _ := user.Current()
	if currentUser != nil {
		kubeConfigPath = path.Join(currentUser.HomeDir, ".kube", "config")
	}

	flags.StringVar(&installCmdOptions.Kube.ConfigPath, "kubeconfig", kubeConfigPath, "Path to kubeconfig for retrieve contexts")
}

package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

// AskAboutKubeContext provide ability select specific context if here few declared in one kubeconfig
func AskAboutKubeContext(kubeOptions *entity.Kube) error {
	kubeConfigPath := kubeOptions.ConfigPath
	if kubeOptions.Context == "" {
		contexts, err := kube.GetAllContexts(kubeConfigPath)
		if err != nil {
			return err
		}

		if len(contexts) == 1 {
			kubeOptions.Context = contexts[0]
		} else {
			_, selectedContext := prompt.NewPrompt().Select(contexts, "Select Kubernetes context")
			kubeOptions.Context = selectedContext
		}

	}
	kubeOptions.ClusterName = kubeOptions.Context
	return nil
}

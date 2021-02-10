package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func AskAboutKubeContext(kubeOptions *install.Kube) error {
	kubeConfigPath := kubeOptions.ConfigPath
	if kubeOptions.Context == "" {
		contexts, err := kube.GetAllContexts(kubeConfigPath)
		if err != nil {
			return err
		}

		if len(contexts) == 1 {
			kubeOptions.Context = contexts[0]
		} else {
			_, selectedContext := prompt.Select(contexts, "Select Kubernetes context")
			kubeOptions.Context = selectedContext
		}

	}
	kubeOptions.ClusterName = kubeOptions.Context
	return nil
}

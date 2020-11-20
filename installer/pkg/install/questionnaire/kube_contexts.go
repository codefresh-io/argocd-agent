package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func AskAboutKubeContext(installOptions *install.InstallCmdOptions) error {
	kubeOptions := installOptions.Kube
	kubeConfigPath := installOptions.Kube.ConfigPath
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
	installOptions.Kube.ClusterName = kubeOptions.Context
	return nil
}

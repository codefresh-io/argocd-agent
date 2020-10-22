package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func AskAboutKubeContext(installOptions *install.InstallCmdOptions) (error, string) {
	kubeOptions := installOptions.Kube
	kubeConfigPath := installOptions.Kube.ConfigPath
	if kubeOptions.Context == "" {
		contexts, err := kube.GetAllContexts(kubeConfigPath)
		if err != nil {
			return err, kubeOptions.Context
		}

		err, selectedContext := prompt.Select(contexts, "Select Kubernetes context")
		kubeOptions.Context = selectedContext
	}
	return nil, kubeOptions.Context
}

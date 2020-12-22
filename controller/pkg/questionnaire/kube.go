package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/controller/pkg/install"
	"github.com/codefresh-io/argocd-listener/controller/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func AskAboutKubeContext(installOptions *install.CmdOptions) error {
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
	installOptions.Kube.Context = kubeOptions.Context
	return nil
}

func AskAboutManifest(installOptions *install.CmdOptions) error {
	return prompt.InputWithDefault(&installOptions.Kube.ManifestPath, "Install manifest path/url", "https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml")
}

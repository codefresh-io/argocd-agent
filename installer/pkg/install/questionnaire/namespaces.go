package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func AskAboutNamespace(installOptions *install.InstallCmdOptions, kubeClient kube.Kube) error {
	namespaces, err := kubeClient.GetNamespaces()
	if err != nil {
		err = prompt.InputWithDefault(&installOptions.Kube.Namespace, "Kubernetes namespace to install", "default")
		if err != nil {
			return err
		}
	} else {
		err, selectedNamespace := prompt.Select(namespaces, "Select Kubernetes namespace")
		if err != nil {
			return err
		}
		installOptions.Kube.Namespace = selectedNamespace
	}
	return nil
}

package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/controller/pkg/install"
	"github.com/codefresh-io/argocd-listener/controller/pkg/kube"
	"github.com/codefresh-io/argocd-listener/controller/pkg/prompt"
)

func AskAboutNamespace(installOptions *install.CmdOptions, kubeClient kube.Kube) error {
	if installOptions.Kube.Namespace == "" {
		namespaces, err := kubeClient.GetNamespaces()
		if err != nil {
			err = prompt.InputWithDefault(&installOptions.Kube.Namespace, "Kubernetes namespace to install", "default")
			if err != nil {
				return err
			}
		} else {
			err, selectedNamespace := prompt.Select(namespaces, "Select the namespace")
			if err != nil {
				return err
			}
			installOptions.Kube.Namespace = selectedNamespace
		}
	}
	return nil
}

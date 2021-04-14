package watch

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/kube"
	"k8s.io/client-go/dynamic"
)

type Watcher interface {
	Watch() error
}

func getKubeconfig() (dynamic.Interface, error) {
	config, err := kube.BuildConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

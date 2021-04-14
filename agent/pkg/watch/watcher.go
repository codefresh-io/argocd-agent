package watch

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/kube"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"time"
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

func start(kubeInformerFactory dynamicinformer.DynamicSharedInformerFactory) {
	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)
}

func getInformer(crd schema.GroupVersionResource) (cache.SharedIndexInformer, dynamicinformer.DynamicSharedInformerFactory, error) {
	clientset, err := getKubeconfig()
	if err != nil {
		return nil, nil, err
	}
	kubeInformerFactory := dynamicinformer.NewDynamicSharedInformerFactory(clientset, time.Minute*30)
	informer := kubeInformerFactory.ForResource(crd).Informer()
	return informer, kubeInformerFactory, nil
}

func Start() error {
	projectWatcher := NewProjectWatcher()
	applicationWatcher := NewApplicationWatcher()

	err := projectWatcher.Watch()
	if err != nil {
		return err
	}

	err = applicationWatcher.Watch()
	if err != nil {
		return err
	}

	return nil
}

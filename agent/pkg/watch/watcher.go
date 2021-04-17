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
	Watch() (dynamicinformer.DynamicSharedInformerFactory, error)
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

func getInformer(crd schema.GroupVersionResource) (cache.SharedIndexInformer, dynamicinformer.DynamicSharedInformerFactory, error) {
	clientset, err := getKubeconfig()
	if err != nil {
		return nil, nil, err
	}
	kubeInformerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute*30, "argocd", nil)
	informer := kubeInformerFactory.ForResource(crd).Informer()
	return informer, kubeInformerFactory, nil
}

func Start() error {
	projectWatcher, err := NewProjectWatcher()
	if err != nil {
		return err
	}

	applicationWatcher, err := NewApplicationWatcher()
	if err != nil {
		return err
	}

	projectInformerFactory, err := projectWatcher.Watch()
	if err != nil {
		return err
	}

	applicationInformerFactory, err := applicationWatcher.Watch()
	if err != nil {
		return err
	}

	stopApplication := make(chan struct{})
	defer close(stopApplication)
	applicationInformerFactory.Start(stopApplication)

	stopProject := make(chan struct{})
	defer close(stopProject)
	projectInformerFactory.Start(stopProject)

	for {
		time.Sleep(time.Second)
	}
}

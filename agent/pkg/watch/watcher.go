package watch

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/kube"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
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

func GetResourceInterface(crd schema.GroupVersionResource, namespace string) (dynamic.ResourceInterface, error) {
	clientset, err := getKubeconfig()
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		return clientset.Resource(crd), nil
	}
	return clientset.Resource(crd).Namespace(namespace), nil
}

func getInformer(crd schema.GroupVersionResource, namespace string) (cache.SharedIndexInformer, dynamicinformer.DynamicSharedInformerFactory, error) {
	clientset, err := getKubeconfig()
	if err != nil {
		return nil, nil, err
	}

	var kubeInformerFactory dynamicinformer.DynamicSharedInformerFactory
	if namespace != "" {
		kubeInformerFactory = dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute*30, namespace, nil)
	} else {
		kubeInformerFactory = dynamicinformer.NewDynamicSharedInformerFactory(clientset, time.Minute*30)
	}

	informer := kubeInformerFactory.ForResource(crd).Informer()
	return informer, kubeInformerFactory, nil
}

func Start(namespace string, sharding *util.Sharding) error {
	projectWatcher, err := NewProjectWatcher(namespace)
	if err != nil {
		return err
	}

	applicationWatcher, err := NewApplicationWatcher(namespace, sharding)
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
	applicationInformerFactory.Start(stopApplication)

	stopProject := make(chan struct{})
	projectInformerFactory.Start(stopProject)

	return nil
}

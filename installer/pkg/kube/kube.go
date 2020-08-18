package kube

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type (
	Kube interface {
		buildClient() (*kubernetes.Clientset, error)
		GetNamespaces() ([]string, error)
		GetClientSet() *kubernetes.Clientset
	}

	kube struct {
		contextName      string
		namespace        string
		pathToKubeConfig string
		inCluster        bool
		clientSet        *kubernetes.Clientset
	}

	Options struct {
		ContextName      string
		Namespace        string
		PathToKubeConfig string
		InCluster        bool
	}
)

func New(o *Options) (Kube, error) {
	client := &kube{
		contextName:      o.ContextName,
		namespace:        o.Namespace,
		pathToKubeConfig: o.PathToKubeConfig,
		inCluster:        o.InCluster,
	}
	clientSet, err := client.buildClient()

	if err != nil {
		return nil, err
	}

	client.clientSet = clientSet

	return client, nil
}

func (k *kube) buildClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if k.inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: k.pathToKubeConfig},
			&clientcmd.ConfigOverrides{
				CurrentContext: k.contextName,
				Context: clientcmdapi.Context{
					Namespace: k.namespace,
				},
			}).ClientConfig()
	}

	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (k *kube) GetNamespaces() ([]string, error) {
	namespaces, err := k.clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	var result []string

	for _, value := range namespaces.Items {
		if value.Name == "default" {
			result = append([]string{"default"}, result...)
			continue
		}
		result = append(result, value.Name)
	}

	return result, nil
}

func (k *kube) GetClientSet() *kubernetes.Clientset {
	return k.clientSet
}

func GetAllContexts(pathToKubeConfig string) ([]string, error) {
	var result []string
	k8scmd := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: pathToKubeConfig},
		&clientcmd.ConfigOverrides{})

	config, err := k8scmd.RawConfig()

	if err != nil {
		return result, err
	}

	if config.CurrentContext != "" {
		result = append(result, config.CurrentContext)
	}

	for k, _ := range config.Contexts {
		if k == config.CurrentContext {
			continue
		}

		result = append(result, k)
	}

	return result, err
}

package kube

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strconv"
)

func BuildConfig() (*rest.Config, error) {
	inCluster, _ := strconv.ParseBool(os.Getenv("IN_CLUSTER"))
	if inCluster {
		return rest.InClusterConfig()
	}

	if os.Getenv("MASTERURL") != "" {
		cfg, err := clientcmd.BuildConfigFromFlags(os.Getenv("MASTERURL"), "")
		if err != nil {
			return nil, err
		}
		cfg.BearerToken = os.Getenv("BEARERTOKEN")
		cfg.Insecure = true
		return cfg, nil
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(
			os.Getenv("HOME"), ".kube", "kubectx",
		)
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

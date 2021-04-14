package watch

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	projectCRD = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "appprojects",
	}
)

type projectWatcher struct {
	codefreshApi codefresh.CodefreshApi
}

func NewProjectWatcher() Watcher {
	return &projectWatcher{codefreshApi: codefresh.GetInstance()}
}

func (projectWatcher *projectWatcher) Watch() error {
	clientset, err := getKubeconfig()
	if err != nil {
		return err
	}

	return nil
}

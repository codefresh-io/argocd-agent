package watch

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
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
	projectInformer, kubeInformerFactory, err := getInformer(projectCRD)
	if err != nil {
		return err
	}

	projectInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			projects, err := argo.GetInstance().GetProjectsWithCredentialsFromStorage()

			if err != nil {
				logger.GetLogger().Errorf("Failed to get projects, reason: %v", err)
				return
			}

			err = util.ProcessDataWithFilter("projects", nil, projects, nil, func() error {
				projects := service.NewArgoResourceService().AdaptArgoProjects(projects)
				return projectWatcher.codefreshApi.SendResources("projects", projects, len(projects))
			})

			if err != nil {
				logger.GetLogger().Errorf("Failed to send projects to codefresh, reason: %v", err)
			}
		},
		DeleteFunc: func(obj interface{}) {
			projects, err := argo.GetInstance().GetProjectsWithCredentialsFromStorage()

			if err != nil {
				//TODO: add error handling
				return
			}

			err = util.ProcessDataWithFilter("projects", nil, projects, nil, func() error {
				projects := service.NewArgoResourceService().AdaptArgoProjects(projects)
				return projectWatcher.codefreshApi.SendResources("projects", projects, len(projects))
			})
			if err != nil {
				logger.GetLogger().Errorf("Failed to send projects to codefresh, reason: %v", err)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
		},
	})

	start(kubeInformerFactory)

	return nil
}

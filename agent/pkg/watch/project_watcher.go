package watch

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
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
	codefreshApi    codefresh.CodefreshApi
	informer        cache.SharedIndexInformer
	informerFactory dynamicinformer.DynamicSharedInformerFactory
	argoApi         argo.ArgoAPI
}

func NewProjectWatcher() (*projectWatcher, error) {
	informer, informerFactory, err := getInformer(projectCRD)
	if err != nil {
		return nil, err
	}
	return &projectWatcher{codefreshApi: codefresh.GetInstance(), informer: informer,
		informerFactory: informerFactory, argoApi: argo.GetInstance()}, nil
}

func (projectWatcher *projectWatcher) add(obj interface{}) {
	projects, err := projectWatcher.argoApi.GetProjectsWithCredentialsFromStorage()

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
}

func (projectWatcher *projectWatcher) delete(obj interface{}) {
	projects, err := projectWatcher.argoApi.GetProjectsWithCredentialsFromStorage()

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
}

func (projectWatcher *projectWatcher) Watch() (dynamicinformer.DynamicSharedInformerFactory, error) {
	projectWatcher.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			projectWatcher.add(obj)
		},
		DeleteFunc: func(obj interface{}) {
			projectWatcher.delete(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
		},
	})

	return projectWatcher.informerFactory, nil
}

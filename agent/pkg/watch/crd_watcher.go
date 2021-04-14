package watch

import (
	//"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/events"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/kube"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/queue"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"github.com/mitchellh/mapstructure"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/cache"
	"time"
)

var (
	applicationCRD = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}

	projectCRD = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "appprojects",
	}
)

var itemQueue *queue.ItemQueue

func watchApplicationChanges() error {
	config, err := kube.BuildConfig()
	if err != nil {
		return err
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	kubeInformerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute*30, "argocd", nil)
	applicationInformer := kubeInformerFactory.ForResource(applicationCRD).Informer()

	api := codefresh2.GetInstance()

	applicationInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var app argoSdk.ArgoApplication
			err := mapstructure.Decode(obj.(*unstructured.Unstructured).Object, &app)

			if err != nil {
				logger.GetLogger().Errorf("Failed to decode argo application, reason: %v", err)
				return
			}

			itemQueue.Enqueue(obj.(*unstructured.Unstructured))

			applications, err := argo.GetInstance().GetApplicationsWithCredentialsFromStorage()

			if err != nil {
				logger.GetLogger().Errorf("Failed to get applications, reason: %v", err)
				return
			}

			err = util.ProcessDataWithFilter("applications", nil, applications, nil, func() error {
				applications := service.NewArgoResourceService().AdaptArgoApplications(applications)
				return api.SendResources("applications", applications, len(applications))
			})

			if err != nil {
				logger.GetLogger().Errorf("Failed to send applications to codefresh, reason: %v", err)
				return
			}

			logger.GetLogger().Info("Successfully sent applications to codefresh")

			applicationCreatedHandler := events.GetApplicationCreatedHandlerInstance()
			err = applicationCreatedHandler.Handle(app)

			if err != nil {
				logger.GetLogger().Errorf("Failed to handle create application event use handler, reason: %v", err)
			} else {
				logger.GetLogger().Infof("Successfully handle new application \"%v\" ", app.Metadata.Name)
			}
		},
		DeleteFunc: func(obj interface{}) {
			var app argoSdk.ArgoApplication
			err := mapstructure.Decode(obj.(*unstructured.Unstructured).Object, &app)
			if err != nil {
				logger.GetLogger().Errorf("Failed to decode argo application, reason: %v", err)
				return
			}

			applications, err := argo.GetInstance().GetApplicationsWithCredentialsFromStorage()
			if err != nil {
				logger.GetLogger().Errorf("Failed to get applications, reason: %v", err)
				return
			}

			err = util.ProcessDataWithFilter("applications", nil, applications, nil, func() error {
				applications := service.NewArgoResourceService().AdaptArgoApplications(applications)
				return api.SendResources("applications", applications, len(applications))
			})

			if err != nil {
				logger.GetLogger().Errorf("Failed to send applications to codefresh, reason: %v", err)
				return
			}

			applicationRemovedHandler := events.GetApplicationRemovedHandlerInstance()
			err = applicationRemovedHandler.Handle(app)

			if err != nil {
				logger.GetLogger().Errorf("Failed to handle remove application event use handler, reason: %v", err)
			}

			err, _ = service.NewGitopsService().MarkEnvAsRemoved(obj)
			if err != nil {
				logger.GetLogger().Errorf("Failed to update application status as 'Deleted', reason: %v", err)
			}

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			itemQueue.Enqueue(newObj.(*unstructured.Unstructured))
		},
	})

	projectInformer := kubeInformerFactory.ForResource(projectCRD).Informer()

	projectInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			projects, err := argo.GetInstance().GetProjectsWithCredentialsFromStorage()

			if err != nil {
				logger.GetLogger().Errorf("Failed to get projects, reason: %v", err)
				return
			}

			err = util.ProcessDataWithFilter("projects", nil, projects, nil, func() error {
				projects := service.NewArgoResourceService().AdaptArgoProjects(projects)
				return api.SendResources("projects", projects, len(projects))
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
				return api.SendResources("projects", projects, len(projects))
			})
			if err != nil {
				logger.GetLogger().Errorf("Failed to send projects to codefresh, reason: %v", err)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
		},
	})

	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)

	for {
		time.Sleep(time.Second)
	}

}

func Watch() error {
	itemQueue = queue.GetInstance()
	return watchApplicationChanges()
}

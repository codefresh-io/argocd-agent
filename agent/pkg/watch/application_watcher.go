package watch

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/events"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/queue"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"github.com/mitchellh/mapstructure"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

var (
	applicationCRD = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}
)

type applicationWatcher struct {
	codefreshApi    codefresh.CodefreshApi
	itemQueue       *queue.ItemQueue
	informer        cache.SharedIndexInformer
	informerFactory dynamicinformer.DynamicSharedInformerFactory
	argoApi         argo.ArgoAPI
}

func NewApplicationWatcher() (Watcher, error) {
	informer, informerFactory, err := getInformer(applicationCRD)
	if err != nil {
		return nil, err
	}
	return &applicationWatcher{codefreshApi: codefresh.GetInstance(), itemQueue: queue.GetInstance(),
		informer: informer, informerFactory: informerFactory, argoApi: argo.GetInstance()}, nil
}

func (watcher *applicationWatcher) add(obj interface{}) {
	var app argoSdk.ArgoApplication
	err := mapstructure.Decode(obj.(*unstructured.Unstructured).Object, &app)

	if err != nil {
		logger.GetLogger().Errorf("Failed to decode argo application, reason: %v", err)
		return
	}

	watcher.itemQueue.Enqueue(obj.(*unstructured.Unstructured))

	applications, err := watcher.argoApi.GetApplicationsWithCredentialsFromStorage()

	if err != nil {
		logger.GetLogger().Errorf("Failed to get applications, reason: %v", err)
		return
	}

	err = util.ProcessDataWithFilter("applications", nil, applications, nil, func() error {
		applications := service.NewArgoResourceService().AdaptArgoApplications(applications)
		return watcher.codefreshApi.SendResources("applications", applications, len(applications))
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
}

func (watcher *applicationWatcher) delete(obj interface{}) {
	var app argoSdk.ArgoApplication
	err := mapstructure.Decode(obj.(*unstructured.Unstructured).Object, &app)
	if err != nil {
		logger.GetLogger().Errorf("Failed to decode argo application, reason: %v", err)
		return
	}

	applications, err := watcher.argoApi.GetApplicationsWithCredentialsFromStorage()
	if err != nil {
		logger.GetLogger().Errorf("Failed to get applications, reason: %v", err)
		return
	}

	err = util.ProcessDataWithFilter("applications", nil, applications, nil, func() error {
		applications := service.NewArgoResourceService().AdaptArgoApplications(applications)
		return watcher.codefreshApi.SendResources("applications", applications, len(applications))
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
}

func (watcher *applicationWatcher) update(newObj interface{}) {
	watcher.itemQueue.Enqueue(newObj.(*unstructured.Unstructured))
}

func (watcher *applicationWatcher) Watch() error {

	watcher.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			watcher.add(obj)
		},
		DeleteFunc: func(obj interface{}) {
			watcher.delete(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			watcher.update(newObj)
		},
	})

	start(watcher.informerFactory)

	return nil
}

package startup

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/events"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/newrelic"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/queue"
	"github.com/codefresh-io/argocd-listener/agent/pkg/scheduler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"github.com/codefresh-io/argocd-listener/agent/pkg/watch"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Runner struct {
	input        *Input
	codefreshApi codefresh.CodefreshApi
}

func NewRunner(input *Input, codefreshApi codefresh.CodefreshApi) *Runner {
	return &Runner{input, codefreshApi}
}

func (runner *Runner) ensureIntegration() {
	input := runner.input
	codefreshApi := runner.codefreshApi
	serverVersion, err := argo.GetUnauthorizedApiInstance().GetVersion(input.argoHost)
	if err != nil {
		err = codefreshApi.CreateIntegration(input.codefreshIntegrationName, input.argoHost,
			input.argoUsername, input.argoPassword, input.argoToken, "",
			"argocd", "")
		if err != nil {
			logger.GetLogger().Errorf("Failed to create integration, reason: %s", err.Error())
		}
		return
	}
	err = codefreshApi.CreateIntegration(input.codefreshIntegrationName, input.argoHost,
		input.argoUsername, input.argoPassword, input.argoToken, serverVersion,
		"argocd", "")
	if err != nil {
		logger.GetLogger().Errorf("Failed to create integration, reason: %s", err.Error())
	}
}

func (runner *Runner) Run() error {
	if runner.input.createIntegrationIfNotExist {
		logger.GetLogger().Info("Create new integration because flag createIntegrationIfNotExist is true")
		go runner.ensureIntegration()
	}

	scheduler.StartHeartBeat()
	go scheduler.GetEnvInitializerScheduler().Run()
	go scheduler.GetIntegrationUpdatedScheduler().Run()
	go scheduler.GetRegenerateTokenScheduler().Run()

	err := events.GetSyncHandlerInstance(codefresh.GetInstance(), argo.GetInstance()).Handle()
	if err != nil {
		logger.GetLogger().Errorf("Failed to run sync handler, reason %v", err)
	}

	err = newrelic.GetInstance().Init(runner.input.newRelicLicense, runner.input.envName)
	if err != nil {
		logger.GetLogger().Errorf("Initialize newrelic error %s", err.Error())
	}

	queueProcessor := queue.EnvQueueProcessor{}
	go queueProcessor.Run()

	applicationCRD := schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}

	resourceInterface, err := watch.GetResourceInterface(applicationCRD, runner.input.namespace)
	if err != nil {
		logger.GetLogger().Errorf("Initialize resource interface failed error %s", err.Error())
	}

	list, err := resourceInterface.List(v1.ListOptions{})
	if err != nil {
		logger.GetLogger().Errorf("failed to get list of applications, reason: %s", err.Error())
	}

	sharding := util.NewSharding(runner.input.replicas)
	sharding.InitApplications(list.Items)

	return watch.Start(runner.input.namespace, sharding)
}

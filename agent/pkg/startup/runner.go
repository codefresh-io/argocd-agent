package startup

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/events"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/queue"
	"github.com/codefresh-io/argocd-listener/agent/pkg/scheduler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/watch"
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
	serverVersion, err := argo.GetInstance().GetVersion()
	if err != nil {
		_ = codefreshApi.CreateIntegration(input.codefreshIntegrationName, input.argoHost,
			input.argoUsername, input.argoPassword, input.argoToken, "",
			"argocd", "")
		return
	}
	_ = codefreshApi.CreateIntegration(input.codefreshIntegrationName, input.argoHost,
		input.argoUsername, input.argoPassword, input.argoToken, serverVersion,
		"argocd", "")
}

func (runner *Runner) Run() error {
	if runner.input.createIntegrationIfNotExist {
		logger.GetLogger().Info("Create new integration because flag createIntegrationIfNotExist is true")
		go runner.ensureIntegration()
	}

	scheduler.StartHeartBeat()
	scheduler.StartEnvInitializer()
	scheduler.StartUpdateIntegration()

	err := events.GetSyncHandlerInstance(codefresh.GetInstance(), argo.GetInstance()).Handle()
	if err != nil {
		logger.GetLogger().Errorf("Failed to run sync handler, reason %v", err)
	}

	queueProcessor := queue.EnvQueueProcessor{}
	go queueProcessor.Run()

	return watch.Watch()
}

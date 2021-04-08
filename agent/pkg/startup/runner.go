package startup

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/event_handler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/queue"
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

	err := event_handler.GetSyncHandlerInstance(codefresh.GetInstance(), argo.GetInstance()).Handle()
	if err != nil {
		logger.GetLogger().Errorf("Failed to run sync handler, reason %v", err)
	}

	queueProcessor := queue.EnvQueueProcessor{}
	go queueProcessor.Run()

	return watch.Watch()
}

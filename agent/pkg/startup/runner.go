package startup

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/event_handler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/queue"
	"github.com/codefresh-io/argocd-listener/agent/pkg/scheduler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/watch"
)

type Runner struct {
	input *Input
}

func NewRunner(input *Input) *Runner {
	return &Runner{input}
}

func (runner *Runner) ensureIntegration() {
	input := runner.input
	serverVersion, err := argo.GetInstance().GetVersion()
	if err != nil {
		return
	}
	err = codefresh.GetInstance().CreateIntegration(input.codefreshIntegrationName, input.argoHost,
		input.argoUsername, input.argoPassword, input.argoToken, serverVersion,
		"argocd", "")
	fmt.Println(err)
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

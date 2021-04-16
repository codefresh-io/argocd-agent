package queue

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/events"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	env2 "github.com/codefresh-io/argocd-listener/agent/pkg/transform/env"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util/comparator"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"time"
)

type QueueProcessor interface {
	Run()
}

type EnvQueueProcessor struct {
}

var envQueueProcessor *EnvQueueProcessor

func (processor *EnvQueueProcessor) New() QueueProcessor {
	if envQueueProcessor == nil {
		envQueueProcessor = &EnvQueueProcessor{}
	}
	return envQueueProcessor
}

func updateEnv(obj *argoSdk.ArgoApplication) (error, *codefreshSdk.Environment) {
	envTransformer := env2.GetEnvTransformerInstance(argo.GetInstance())
	err, env := envTransformer.PrepareEnvironment(*obj)
	if err != nil {
		return err, env
	}

	envComparator := comparator.EnvComparator{}

	err = util.ProcessDataWithFilter("environment", &env.Name, env, envComparator.Compare, func() error {
		return events.GetRolloutEventHandlerInstance().Handle(env)
	})

	return err, env
}

func (processor *EnvQueueProcessor) Run() {
	itemQueue := GetInstance()
	for true {
		if itemQueue.Size() > 0 {
			item := itemQueue.Dequeue()
			err, _ := updateEnv(item)
			if err != nil {
				logger.GetLogger().Errorf("Failed to update environment, reason: %v", err)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

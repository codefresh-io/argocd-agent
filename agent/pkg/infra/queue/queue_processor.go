package queue

import (
	"time"

	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/events"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util/comparator"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type QueueProcessor interface {
	Run()
}

type EnvQueueProcessor struct {
	argoApi argo.ArgoAPI
}

var envQueueProcessor *EnvQueueProcessor

func (processor *EnvQueueProcessor) New() QueueProcessor {
	if envQueueProcessor == nil {
		envQueueProcessor = &EnvQueueProcessor{
			argoApi: argo.GetInstance(),
		}
	}
	return envQueueProcessor
}

func updateEnv(obj *argoSdk.ArgoApplication, historyId int64, argoApi argo.ArgoAPI) (error, *codefreshSdk.Environment) {
	envTransformer := service.GetEnvTransformerInstance(argoApi)
	err, envWrapper := envTransformer.PrepareEnvironment(*obj, historyId)
	if err != nil {
		return err, nil
	}

	env := &envWrapper.Environment

	envComparator := comparator.EnvComparator{}

	err = util.ProcessDataWithFilter("environment", &env.Name, env, envComparator.Compare, func() error {
		return events.GetRolloutEventHandlerInstance().Handle(envWrapper)
	})

	return err, env
}

func (processor *EnvQueueProcessor) Run() {
	itemQueue := GetInstance()
	for {
		processStartTime := time.Now()

		if itemQueue.Size() > 0 {
			item := itemQueue.Dequeue()

			dequeueTime := time.Since(processStartTime)
			logger.GetLogger().Debugf("[queue_processor_metric] dequeued in %s", dequeueTime)

			if item != nil {
				logger.GetLogger().Debugf("[queue_processor_metric] processing application %v", item.Application.Metadata.Name)

				err, _ := updateEnv(&item.Application, item.HistoryId, processor.argoApi)

				updateTime := time.Since(processStartTime.Add(dequeueTime))
				logger.GetLogger().Debugf("[queue_processor_metric] env updated in %s", updateTime)

				if err != nil {
					logger.GetLogger().Errorf("Failed to update environment, reason: %v", err)
				}

			}
			logger.GetLogger().Debugf("[queue_processor_metric] application processed in %s in total", time.Since(processStartTime))

			logger.GetLogger().Infof("Queue size %v", itemQueue.Size())
			// don't sleep at all in case there are more items in queue
		} else {
			logger.GetLogger().Debug("queue is empty, standby for 1 second...")
			// sleep in case queue empty
			time.Sleep(1 * time.Second)
		}
	}
}

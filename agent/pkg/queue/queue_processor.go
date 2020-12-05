package queue

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/transform"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util/comparator"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func updateEnv(obj *unstructured.Unstructured) (error, *codefresh.Environment) {
	envTransformer := transform.GetEnvTransformerInstance(argo.GetInstance())
	err, env := envTransformer.PrepareEnvironment(obj.Object)
	if err != nil {
		return err, env
	}

	envComparator := comparator.EnvComparator{}

	err = util.ProcessDataWithFilter("environment", &env.Name, env, envComparator.Compare, func() error {
		_, err = codefresh.GetInstance().SendEnvironment(*env)
		return err
	})

	return nil, env
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

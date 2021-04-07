package main

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/event_handler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/heartbeat"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/queue"
	"github.com/codefresh-io/argocd-listener/agent/pkg/scheduler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/startup"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/watch"
)

func ensureIntegration(integration, host, username, password, token, provider, clusterName string) {
	serverVersion, err := argo.GetInstance().GetVersion()
	if err != nil {
		return
	}
	err = codefresh2.GetInstance().CreateIntegration(integration, host,
		username, password, token, serverVersion,
		provider, clusterName)
	fmt.Println(err)
}

func handleError(err error) {
	store.SetHeartbeatError(err.Error())
	heartbeat.HeartBeatTask()
	// send heartbeat to codefresh before die
	panic(err)
}

func main() {

	input := startup.NewInputFactory().Create()

	err := startup.NewInputValidator(input).Validate()
	if err != nil {
		handleError(err)
	}

	err = startup.NewInputStore(input).Store()
	if err != nil {
		handleError(err)
	}

	//	ensureIntegration(codefreshIntegrationName, argoHost, argoUsername, argoPassword, argoToken, "argocd", "");

	scheduler.StartHeartBeat()
	scheduler.StartEnvInitializer()
	scheduler.StartUpdateIntegration()

	err = event_handler.GetSyncHandlerInstance(codefresh2.GetInstance(), argo.GetInstance()).Handle()
	if err != nil {
		logger.GetLogger().Errorf("Failed to run sync handler, reason %v", err)
	}

	queueProcessor := queue.EnvQueueProcessor{}
	go queueProcessor.Run()

	err = watch.Watch()
	if err != nil {
		logger.GetLogger().Errorf("Cant run agent because %v", err.Error())
		store.SetHeartbeatError(err.Error())
		heartbeat.HeartBeatTask()
		panic(err)
	}

}

package main

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/heartbeat"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/startup"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
)

func handleError(err error) {
	logger.GetLogger().Errorf("Cant run agent because %v", err.Error())
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

	err = startup.NewRunner(input).Run()
	if err != nil {
		handleError(err)
	}

}

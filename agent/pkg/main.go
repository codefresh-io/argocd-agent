package main

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/codefresh-io/argocd-listener/agent/pkg/startup"
)

func handleError(err error) {
	logger.GetLogger().Errorf("Cant run agent because %v", err.Error())
	store.SetHeartbeatError(err.Error())
	service.HeartBeatTask()
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

	err = startup.NewRunner(input, codefresh.GetInstance()).Run()
	if err != nil {
		handleError(err)
	}

}

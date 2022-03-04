package main

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/codefresh-io/argocd-listener/agent/pkg/startup"
	"net/http"
	_ "net/http/pprof"
)

func handleError(err error) {
	hb := service.New()

	logger.GetLogger().Errorf("Cant run agent because %v", err.Error())
	store.SetHeartbeatError(err.Error())

	hb.HeartBeatTask()
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

	http.ListenAndServe("0.0.0.0:6060", nil)

}

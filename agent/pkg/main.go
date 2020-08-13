package main

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/extract"
	"github.com/codefresh-io/argocd-listener/agent/pkg/heartbeat"
	"github.com/codefresh-io/argocd-listener/agent/pkg/scheduler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"os"
)

func main() {

	argoHost, argoHostExistence := os.LookupEnv("ARGO_HOST")
	if !argoHostExistence {
		panic(errors.New("ARGO_HOST variable doesnt exist"))
	}

	argoUsername, argoUsernameExistence := os.LookupEnv("ARGO_USERNAME")
	if !argoUsernameExistence {
		panic(errors.New("ARGO_USERNAME variable doesnt exist"))
	}

	argoPassword, argoPasswordExistence := os.LookupEnv("ARGO_PASSWORD")
	if !argoPasswordExistence {
		panic(errors.New("ARGO_PASSWORD variable doesnt exist"))
	}

	codefreshToken, codefreshTokenExistence := os.LookupEnv("CODEFRESH_TOKEN")
	if !codefreshTokenExistence {
		panic(errors.New("CODEFRESH_TOKEN variable doesnt exist"))
	}

	codefreshHost, codefreshHostExistance := os.LookupEnv("CODEFRESH_HOST")
	if !codefreshHostExistance {
		codefreshHost = "https://g.codefresh.io"
	}

	codefreshIntegrationName, codefreshIntegrationNameExistence := os.LookupEnv("CODEFRESH_INTEGRATION")
	if !codefreshIntegrationNameExistence {
		panic(errors.New("CODEFRESH_INTEGRATION variable doesnt exist"))
	}

	store.SetCodefresh(codefreshHost, codefreshToken, codefreshIntegrationName)

	token, err := argo.GetToken(argoUsername, argoPassword, argoHost)
	if err != nil {
		store.SetHeartbeatError(err.Error())
		heartbeat.HeartBeatTask()
		// send heartbeat to codefresh before die
		panic(err)
	}

	store.SetArgo(token, argoHost)

	scheduler.StartHeartBeat()

	extract.Watch()

	//go argo.Schedule()

	//time.Sleep(100 * time.Minute)
}

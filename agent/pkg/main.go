package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/extract"
	"github.com/codefresh-io/argocd-listener/agent/pkg/handler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/heartbeat"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/scheduler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"os"
)

func main() {

	argoHost, argoHostExistence := os.LookupEnv("ARGO_HOST")
	if !argoHostExistence {
		panic(errors.New("ARGO_HOST variable doesnt exist"))
	}

	argoToken, argoTokenExistence := os.LookupEnv("ARGO_TOKEN")
	if !argoTokenExistence || argoToken == "" {

		argoUsername, argoUsernameExistence := os.LookupEnv("ARGO_USERNAME")
		if !argoUsernameExistence {
			panic(errors.New("ARGO_USERNAME variable doesnt exist"))
		}

		argoPassword, argoPasswordExistence := os.LookupEnv("ARGO_PASSWORD")
		if !argoPasswordExistence {
			panic(errors.New("ARGO_PASSWORD variable doesnt exist"))
		}

		token, err := argo.GetToken(argoUsername, argoPassword, argoHost)

		if err != nil {
			store.SetHeartbeatError(err.Error())
			heartbeat.HeartBeatTask()
			// send heartbeat to codefresh before die
			panic(err)
		}

		store.SetArgo(token, argoHost)

	} else {
		store.SetArgo(argoToken, argoHost)
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

	var applications []string
	syncMode, _ := os.LookupEnv("SYNC_MODE")
	if syncMode == codefresh2.SelectSync {
		applicationsToSyncEncodedJson, _ := os.LookupEnv("APPLICATIONS_FOR_SYNC")
		applicationsToSyncJson, _ := base64.StdEncoding.DecodeString(applicationsToSyncEncodedJson)
		_ = json.Unmarshal(applicationsToSyncJson, &applications)
	}

	store.SetSyncOptions(syncMode, applications)

	store.SetCodefresh(codefreshHost, codefreshToken, codefreshIntegrationName)

	agentVersion, agentVersionExistence := os.LookupEnv("AGENT_VERSION")
	if !agentVersionExistence {
		logger.GetLogger().Errorf("No agent version!")
	} else {
		store.SetAgent(agentVersion)
	}

	//  @todo - move codefresh git integration token to env during installation
	err, contextPayload := codefresh2.GetInstance().GetDefaultGitContext()
	if err != nil {
		logger.GetLogger().Errorf("Failed to get git context, reason: %v", err)
	} else {
		store.SetGit(contextPayload.Spec.Data.Auth.Password)
	}

	scheduler.StartHeartBeat()
	scheduler.StartEnvInitializer()

	err = handler.GetSyncHandlerInstance(codefresh2.GetInstance(), argo.GetInstance()).Handle()
	if err != nil {
		logger.GetLogger().Errorf("Failed to run sync handler, reason %v", err)
	}

	err = extract.Watch()
	if err != nil {
		logger.GetLogger().Errorf("Cant run agent because %v", err.Error())
		store.SetHeartbeatError(err.Error())
		heartbeat.HeartBeatTask()
		panic(err)
	}

}

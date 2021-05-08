package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
)

type regenerateTokenScheduler struct {
	argoApi argo.UnauthorizedApi
}

func GetRegenerateTokenScheduler() Scheduler {
	return &regenerateTokenScheduler{
		argoApi: argo.GetUnauthorizedApiInstance(),
	}
}

func (tokenScheduler *regenerateTokenScheduler) getTime() string {
	return "@every 30m"
}

func (tokenScheduler *regenerateTokenScheduler) getFunc() func() {
	return tokenScheduler.regenerateToken
}

//Run start lister about new environments and update their statuses
func (tokenScheduler *regenerateTokenScheduler) Run() {
	run(tokenScheduler)
}

func (tokenScheduler *regenerateTokenScheduler) regenerateToken() {
	argoValues := store.GetStore().Argo
	if argoValues.Username != "" || argoValues.Password != "" {
		logger.GetLogger().Info("Regenerate argocd token")
		token, err := tokenScheduler.argoApi.GetToken(argoValues.Username, argoValues.Password, argoValues.Host)
		if err != nil {
			store.SetArgoToken(token)
		}
	}
}

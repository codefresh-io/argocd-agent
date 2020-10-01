package heartbeat

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
)

func HeartBeatTask() {
	err := codefresh.GetInstance().HeartBeat(store.GetStore().Heartbeat.Error)
	if err != nil {
		logger.GetLogger().Errorf("Failed to send heartbeat status, reason %v", err)
	}
}

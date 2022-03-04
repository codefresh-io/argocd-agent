package service

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
)

var heartbeatAmount = 0

type heartbeat struct {
	Api *codefresh.Api
}

func New() *heartbeat {
	return &heartbeat{
		Api: codefresh.GetInstance(),
	}
}

func (h *heartbeat) HeartBeatTask() {
	err := h.Api.HeartBeat(store.GetStore().Heartbeat.Error)
	if err != nil {
		logger.GetLogger().Errorf("Failed to send heartbeat status, reason %v", err)
	}

	heartbeatAmount++
	if heartbeatAmount%100 == 0 {
		logger.GetLogger().Infof("Im still alive, heartbeat amount %v", heartbeatAmount)
	}
}

package heartbeat

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
)

func HeartBeatTask() {
	err := codefresh.GetInstance().HeartBeat(store.GetStore().Heartbeat.Error)
	if err != nil {
		fmt.Println(err)
	}
}

package argo

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"time"
)

var HeartBeatInterval = 5 * time.Second

func StartHeartBeat() {
	ticker := time.NewTicker(HeartBeatInterval)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				err := codefresh.GetInstance().HeartBeat()
				if err != nil {
					fmt.Println(err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

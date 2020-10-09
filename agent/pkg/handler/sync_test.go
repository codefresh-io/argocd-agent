package handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"testing"
)

func TestSyncWithNoneMode(t *testing.T) {

	store.SetSyncOptions(codefresh.None, []string{})

	err := syncHandler.Handle()
	if err != nil {
		t.Error(err)
	}

}

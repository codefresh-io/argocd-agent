package handler

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/uninstall"
	"testing"
)

func TestRun(t *testing.T) {
	client := New(uninstall.CmdOptions{})
	err := client.Run()
	if err != nil && err.Error() != "please provide options to select from" {
		t.Errorf("'TestRun' failed, unexpected error, expected '%v', got '%v'", "please provide options to select from", err.Error())
	}

}

package store

import (
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

func TestSetGitContext(t *testing.T) {

	SetGitContext(codefresh.ContextPayload{
		Metadata: struct {
			Name string `json:"name"`
		}{},
		Spec: struct {
			Type string `json:"type"`
			Data struct {
				Auth struct {
					Type           string `json:"type"`
					Username       string `json:"username"`
					Password       string `json:"password"`
					ApiHost        string `json:"apiHost"`
					ApiPathPrefix  string `json:"apiPathPrefix"`
					SshPrivateKey  string `json:"sshPrivateKey"`
					AppId          string `json:"appId"`
					InstallationId string `json:"installationId"`
					PrivateKey     string `json:"privateKey"`
				} `json:"auth"`
			} `json:"data"`
		}{},
	})

	ctx := &GetStore().Git.Context
	if ctx == nil {
		t.Error("Context should exist")
	}
}

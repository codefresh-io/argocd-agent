package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestAskAboutArgoCredentials(t *testing.T) {

	installCmdOptions := &install.InstallCmdOptions{
		Argo: struct {
			Host     string
			Username string
			Password string
			Token    string
			Update   bool
		}{Host: "https://localhost", Username: "test", Password: "test", Token: "test", Update: false},
	}

	_ = AskAboutArgoCredentials(installCmdOptions, nil)

	if installCmdOptions.Argo.Host != "https://localhost" {
		t.Errorf("Argocd host shouldnt be changed in case if it is passed from cli")
	}
}

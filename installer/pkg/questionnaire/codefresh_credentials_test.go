package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/cliconfig"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"testing"
)

type MockCliConfig struct {
}

func (cf *MockCliConfig) GetCurrentConfig() (*cliconfig.CliConfigItem, error) {

	return &cliconfig.CliConfigItem{
		Name:  "test",
		Url:   "http://localhost",
		Token: "token",
	}, nil
}

func TestAskAboutCodefreshCredentialsWithoutPredefinedCreds(t *testing.T) {
	questionnaire := CodefreshCredentialsQuestionnaire{cliConfig: &MockCliConfig{}}
	installOptions := &entity.InstallCmdOptions{}
	err := questionnaire.AskAboutCodefreshCredentials(installOptions)
	if err != nil {
		t.Error("Cant generate codefresh credentials")
	}
	if installOptions.Codefresh.Host != "http://localhost" {
		t.Error("Wrong codefresh host")
	}
	if installOptions.Codefresh.Token != "token" {
		t.Error("Wrong codefresh token")
	}
}

func TestAskAboutCodefreshCredentialsWithPredefinedCreds(t *testing.T) {
	questionnaire := CodefreshCredentialsQuestionnaire{cliConfig: &MockCliConfig{}}
	installOptions := &entity.InstallCmdOptions{}
	installOptions.Codefresh.Token = "token"
	installOptions.Codefresh.Host = "http://localhost"
	err := questionnaire.AskAboutCodefreshCredentials(installOptions)
	if err != nil {
		t.Error("Cant generate codefresh credentials")
	}
	if installOptions.Codefresh.Host != "http://localhost" {
		t.Error("Wrong codefresh host")
	}
	if installOptions.Codefresh.Token != "token" {
		t.Error("Wrong codefresh token")
	}
}

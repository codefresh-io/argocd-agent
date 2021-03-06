package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/cliconfig"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
)

type CodefreshCredentialsQuestionnaire struct {
	cliConfig cliconfig.CliConfig
}

func NewCodefreshCredentialsQuestionnaire() *CodefreshCredentialsQuestionnaire {
	return &CodefreshCredentialsQuestionnaire{cliConfig: cliconfig.NewCliConfig()}
}

// AskAboutCodefreshCredentials request codefresh credentials if it wasnt passed in cli , by default we are taking it from codefresh config file
func (credentialsQuestionnaire *CodefreshCredentialsQuestionnaire) AskAboutCodefreshCredentials(installOptions *entity.InstallCmdOptions) error {
	if installOptions.Codefresh.Token == "" || installOptions.Codefresh.Host == "" {
		config, err := credentialsQuestionnaire.cliConfig.GetCurrentConfig()
		if err != nil {
			return err
		}
		installOptions.Codefresh.Token = config.Token
		installOptions.Codefresh.Host = config.Url
	}
	return nil
}

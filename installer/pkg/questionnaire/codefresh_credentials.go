package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/cliconfig"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/type"
)

// AskAboutCodefreshCredentials request codefresh credentials if it wasnt passed in cli , by default we are taking it from codefresh config file
func AskAboutCodefreshCredentials(installOptions *_type.InstallCmdOptions) error {
	if installOptions.Codefresh.Token == "" || installOptions.Codefresh.Host == "" {
		config, err := cliconfig.GetCurrentConfig()
		if err != nil {
			return err
		}
		installOptions.Codefresh.Token = config.Token
		installOptions.Codefresh.Host = config.Url
	}
	return nil
}

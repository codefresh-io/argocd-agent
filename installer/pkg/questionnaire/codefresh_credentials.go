package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/cliconfig"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
)

func AskAboutCodefreshCredentials(installOptions *install.InstallCmdOptions) error {
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

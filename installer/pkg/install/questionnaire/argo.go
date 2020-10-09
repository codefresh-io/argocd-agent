package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"regexp"
)

func AskAboutArgoCredentials(installOptions *install.InstallCmdOptions) error {
	err := prompt.Input(&installOptions.Argo.Host, "Argo host, example: https://example.com")
	if err != nil {
		return err
	}

	withProtocol, err := regexp.MatchString("^https?://", installOptions.Argo.Host)
	if err != nil {
		return err
	}

	// customer not put protocol during installation
	if !withProtocol {
		installOptions.Argo.Host = "https://" + installOptions.Argo.Host
	}

	// removing / in the end
	installOptions.Argo.Host = regexp.MustCompile("/+$").ReplaceAllString(installOptions.Argo.Host, "")

	err, useArgocdToken := prompt.Confirm("Do you want use argocd auth token instead username/password auth?")
	if err != nil {
		return err
	}

	if useArgocdToken {
		err = prompt.InputWithDefault(&installOptions.Argo.Token, "Argo token", "")
		if err != nil {
			return err
		}
	} else {
		err = prompt.InputWithDefault(&installOptions.Argo.Username, "Argo username", "admin")
		if err != nil {
			return err
		}

		err = prompt.InputPassword(&installOptions.Argo.Password, "Argo password")
		if err != nil {
			return err
		}
	}

	return nil
}

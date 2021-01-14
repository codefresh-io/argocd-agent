package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"regexp"
)

func AskAboutArgoCredentials(installOptions *install.InstallCmdOptions) error {

	err := prompt.Input(&installOptions.Argo.Host, "Argo host, for example: https://example.com")
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

	if (installOptions.Argo.Token != "") || ((installOptions.Argo.Username != "") && (installOptions.Argo.Password != "")) {
		return nil
	}

	//err, useArgocdToken := prompt.Confirm("Choose an authentication method")
	useArgocdToken := "Auth token - Recommended [https://codefresh.io/docs/docs/ci-cd-guides/gitops-deployments/]"
	useUserAndPass := "Username and password"
	authenticationMethodOptions := []string{useArgocdToken, useUserAndPass}
	err, authenticationMethod := prompt.Select(authenticationMethodOptions, "Choose an authentication method")
	if err != nil {
		return err
	}

	if authenticationMethod == useArgocdToken {
		err = prompt.InputWithDefault(&installOptions.Argo.Token, "Argo token", "")
		if err != nil {
			return err
		}
	} else if authenticationMethod == useUserAndPass {
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

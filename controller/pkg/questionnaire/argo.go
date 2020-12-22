package questionnaire

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/controller/pkg/install"
	"github.com/codefresh-io/argocd-listener/controller/pkg/logger"
	"github.com/codefresh-io/argocd-listener/controller/pkg/prompt"
)

func AskAboutPass(installOptions *install.CmdOptions) error {
	i := 1
	var MaxRetry = 3
	if installOptions.Argo.Password != "" {
		return nil
	}

	for i < MaxRetry {
		i++
		installOptions.Argo.Password = askAboutPass()
		if installOptions.Argo.Password != "" {
			return nil
		}
	}

	return errors.New("passwords are different")
}

func askAboutPass() string {
	var firstPassword string
	var secondPassword string

	_ = prompt.InputPassword(&firstPassword, "Argo password")
	_ = prompt.InputPassword(&secondPassword, "Confirm password")
	if firstPassword != secondPassword {
		logger.Error("Passwords are different")
		return ""
	} else if firstPassword == "" {
		logger.Error("Passwords is too short")
		return ""
	}
	return firstPassword
}

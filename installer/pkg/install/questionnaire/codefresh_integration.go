package questionnaire

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func ensureIntegration(installCmdOptions *install.InstallCmdOptions) error {
	serverVersion, err := argo.GetInstance().GetVersion()
	if err != nil {
		return err
	}

	err = holder.ApiHolder.CreateIntegration(installCmdOptions.Codefresh.Integration, installCmdOptions.Argo.Host, installCmdOptions.Argo.Username, installCmdOptions.Argo.Password, installCmdOptions.Argo.Token, serverVersion)
	if err == nil {
		return nil
	}

	codefreshErr, ok := err.(*codefresh.CodefreshError)
	if !ok {
		return err
	}

	if codefreshErr.Status != 409 {
		return codefreshErr
	}

	needUpdate := installCmdOptions.Argo.Update
	if !needUpdate {
		err, needUpdate = prompt.Confirm("You already have integration with this name, do you want to update it")
		if err != nil {
			return err
		}
	}

	if !needUpdate {
		return fmt.Errorf("you should update integration")
	}

	err = holder.ApiHolder.UpdateIntegration(installCmdOptions.Codefresh.Integration, installCmdOptions.Argo.Host, installCmdOptions.Argo.Username, installCmdOptions.Argo.Password, installCmdOptions.Argo.Token, serverVersion)

	if err != nil {
		return err
	}

	return nil
}

func AskAboutCodefreshIntegration(installOptions *install.InstallCmdOptions) error {
	if installOptions.Agent.Interactive {
		err := prompt.InputWithDefault(&installOptions.Codefresh.Integration, "Codefresh integration name", "argocd")
		if err != nil {
			return err
		}
	}

	return ensureIntegration(installOptions)
}

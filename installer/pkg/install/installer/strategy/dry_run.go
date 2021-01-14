package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/acceptance_tests"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/helper"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/questionnaire"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/fatih/structs"
)

type DryRunArgoCdAgentInstaller struct {
}

func (installer *DryRunArgoCdAgentInstaller) Install(installCmdOptions install.InstallCmdOptions) (error, string) {
	var err error
	// should be in beg for show correct events
	_ = questionnaire.AskAboutCodefreshCredentials(&installCmdOptions)

	holder.ApiHolder = codefresh.Api{
		Token:       installCmdOptions.Codefresh.Token,
		Host:        installCmdOptions.Codefresh.Host,
		Integration: installCmdOptions.Codefresh.Integration,
	}

	kubeOptions := installCmdOptions.Kube
	if kubeOptions.Namespace == "" {
		kubeOptions.Namespace = "argocd"
	}

	if installCmdOptions.Agent.Interactive {
		err = questionnaire.AskAboutArgoCredentials(&installCmdOptions)
		if err != nil {
			sendArgoAgentInstalledEvent(FAILED, err.Error())
			return err, ""
		}
	}

	err = acceptance_tests.New().VerifyArgoSetup(&installCmdOptions.Argo)
	if err != nil {
		msg := fmt.Sprintf("Testing argo requirements failed - \"%s\"", err.Error())
		sendArgoAgentInstalledEvent(FAILED, msg)
		return errors.New(msg), ""
	}

	err = questionnaire.AskAboutCodefreshIntegration(&installCmdOptions)
	if err != nil {
		sendArgoAgentInstalledEvent(FAILED, err.Error())
		return err, ""
	}

	err = questionnaire.AskAboutGitContext(&installCmdOptions)
	if err != nil {
		sendArgoAgentInstalledEvent(FAILED, err.Error())
		return err, ""
	}

	questionnaire.AskAboutSyncOptions(&installCmdOptions)

	installCmdOptions.Codefresh.Token = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Codefresh.Token))
	installCmdOptions.Argo.Token = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Argo.Token))
	installCmdOptions.Argo.Password = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Argo.Password))

	helper.ShowSummary(&installCmdOptions)

	sendArgoAgentInstalledEvent(SUCCESS, "")

	installOptions := templates.InstallOptions{
		Templates:      kubernetes.TemplatesMap(),
		TemplateValues: structs.Map(installCmdOptions),
		Namespace:      kubeOptions.Namespace,
	}

	err, manifest := templates.GenerateManifest(&installOptions)

	return nil, manifest
}

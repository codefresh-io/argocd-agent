package install

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	cfEventSender "github.com/codefresh-io/argocd-listener/installer/pkg/cfeventsender"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/acceptance"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/helper"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/questionnaire"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/fatih/structs"
)

type InstallIntegration struct {
	argoApi      argo.ArgoAPI
	codefreshApi codefresh.CodefreshApi
}

func Run(installCmdOptions entity.InstallCmdOptions) (error, string) {
	// should be in begining for show correct events
	_ = questionnaire.NewCodefreshCredentialsQuestionnaire().AskAboutCodefreshCredentials(&installCmdOptions)
	store.SetCodefresh(installCmdOptions.Codefresh.Host, installCmdOptions.Codefresh.Token, installCmdOptions.Codefresh.Integration)

	var err error
	eventSender := cfEventSender.New(cfEventSender.EVENT_AGENT_INSTALL)

	kubeConfigPath := installCmdOptions.Kube.ConfigPath
	kubeOptions := installCmdOptions.Kube

	_ = questionnaire.AskAboutKubeContext(&kubeOptions)
	clusterName := kubeOptions.Context
	kubeClient, err := kube.New(&kube.Options{
		ContextName:      kubeOptions.Context,
		Namespace:        kubeOptions.Namespace,
		PathToKubeConfig: kubeConfigPath,
		InCluster:        kubeOptions.InCluster,
	})
	if err != nil {
		return err, ""
	}
	_ = questionnaire.AskAboutNamespace(&installCmdOptions.Kube, kubeClient, true)

	kubeOptions = installCmdOptions.Kube

	err = prompt.NewPrompt().InputWithDefault(&installCmdOptions.Codefresh.Integration, "Codefresh integration name", "argocd")
	if err != nil {
		return err, ""
	}

	err = questionnaire.NewArgoQuestionnaire().AskAboutArgoCredentials(&installCmdOptions, kubeClient)
	if err != nil {
		eventSender.Fail(err.Error())
		return errors.New(err.Error()), ""
	}

	err = acceptance.New().Verify(&installCmdOptions.Argo)
	if err != nil {
		msg := fmt.Sprintf("Testing requirements failed - \"%s\"", err.Error())
		eventSender.Fail(msg)
		return errors.New(msg), ""
	}

	_ = questionnaire.NewGitContextQuestionnaire().AskAboutGitContext(&installCmdOptions)

	// Need check if we want support not in cluster mode with Product owner
	installCmdOptions.Kube.InCluster = true

	questionnaire.AskAboutSyncOptions(&installCmdOptions)

	installIntegration := &InstallIntegration{codefreshApi: codefresh.GetInstance(), argoApi: argo.GetInstance()}

	err = installIntegration.ensureIntegration(&installCmdOptions, clusterName)
	if err != nil {
		eventSender.Fail(err.Error())
		return err, ""
	}

	installCmdOptions.Codefresh.Token = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Codefresh.Token))
	installCmdOptions.Argo.Token = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Argo.Token))
	installCmdOptions.Argo.Password = base64.StdEncoding.EncodeToString([]byte(installCmdOptions.Argo.Password))

	installOptions := templates.InstallOptions{
		Templates:        kubernetes.TemplatesMap(),
		TemplateValues:   structs.Map(installCmdOptions),
		Namespace:        kubeOptions.Namespace,
		KubeClientSet:    kubeClient.GetClientSet(),
		KubeCrdClientSet: kubeClient.GetCrdClientSet(),
		KubeManifestPath: installCmdOptions.Kube.ManifestPath,
	}
	helper.ShowSummary(&installCmdOptions)

	var kind, name, manifest string

	if installOptions.KubeManifestPath != "" {
		err, kind, name, manifest = templates.DryRunInstall(&installOptions)
	} else {
		err, kind, name, manifest = templates.Install(&installOptions)
	}

	if err != nil {
		msg := fmt.Sprintf("Argo agent installation resource \"%s\" with name \"%s\" finished with error , reason: %v ", kind, name, err)
		eventSender.Fail(err.Error())
		return errors.New(msg), ""
	}

	eventSender.Success("Successfully install argocd agent")

	logger.Success(fmt.Sprintf("Argo agent installation finished successfully to namespace \"%s\"", kubeOptions.Namespace))
	logger.Success(fmt.Sprintf("Gitops view: \"%s/gitops\"", installCmdOptions.Codefresh.Host))
	logger.Success(fmt.Sprintf("Documentation: \"%s\"", "https://codefresh.io/docs/docs/ci-cd-guides/gitops-deployments/"))
	return nil, manifest
}

func (installCmd *InstallIntegration) ensureIntegration(installCmdOptions *entity.InstallCmdOptions, clusterName string) error {
	serverVersion, err := installCmd.argoApi.GetVersion()
	if err != nil {
		return err
	}
	err = installCmd.codefreshApi.CreateIntegration(installCmdOptions.Codefresh.Integration, installCmdOptions.Argo.Host,
		installCmdOptions.Argo.Username, installCmdOptions.Argo.Password, installCmdOptions.Argo.Token, serverVersion,
		installCmdOptions.Codefresh.Provider, clusterName)
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
		err, needUpdate = prompt.NewPrompt().Confirm("You already have integration with this name, do you want to update it")
		if err != nil {
			return err
		}
	}

	if !needUpdate {
		return fmt.Errorf("you should update integration")
	}

	err = installCmd.codefreshApi.UpdateIntegration(installCmdOptions.Codefresh.Integration, installCmdOptions.Argo.Host,
		installCmdOptions.Argo.Username, installCmdOptions.Argo.Password, installCmdOptions.Argo.Token, serverVersion,
		installCmdOptions.Codefresh.Provider, clusterName)

	if err != nil {
		return err
	}

	return nil
}

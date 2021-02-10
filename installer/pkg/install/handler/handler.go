package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	cfEventSender "github.com/codefresh-io/argocd-listener/installer/pkg/cf_event_sender"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/acceptance_tests"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/helper"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/questionnaire"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/fatih/structs"
)

func Run(installCmdOptions install.InstallCmdOptions) (error, string) {
	var err error
	eventSender := cfEventSender.New(cfEventSender.EVENT_INSTALL)
	// should be in beg for show correct events
	_ = questionnaire.AskAboutCodefreshCredentials(&installCmdOptions)

	holder.ApiHolder = codefresh.Api{
		Token:       installCmdOptions.Codefresh.Token,
		Host:        installCmdOptions.Codefresh.Host,
		Integration: installCmdOptions.Codefresh.Integration,
	}

	kubeConfigPath := installCmdOptions.Kube.ConfigPath
	kubeOptions := installCmdOptions.Kube

	_ = questionnaire.AskAboutKubeContext(&kubeOptions)

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

	argoServerSvc, err := kubeClient.GetArgoServerSvc(kubeOptions.Namespace)

	if err != nil {
		msg := fmt.Sprintf("We didn't find ArgoCD on \"%s/%s\"", installCmdOptions.Kube.ClusterName, kubeOptions.Namespace)
		eventSender.Fail(msg)
		return errors.New(msg), ""
	} else {
		if kube.IsLoadBalancer(argoServerSvc) {
			balancerHost, _ := kubeClient.GetLoadBalancerHost(argoServerSvc)
			if balancerHost != "" {
				installCmdOptions.Argo.Host = balancerHost
			}
		}
	}

	err = prompt.InputWithDefault(&installCmdOptions.Codefresh.Integration, "Codefresh integration name", "argocd")
	if err != nil {
		return err, ""
	}

	_ = questionnaire.AskAboutArgoCredentials(&installCmdOptions)

	err = acceptance_tests.New().Verify(&installCmdOptions.Argo)
	if err != nil {
		msg := fmt.Sprintf("Testing requirements failed - \"%s\"", err.Error())
		eventSender.Fail(msg)
		return errors.New(msg), ""
	}

	_ = questionnaire.AskAboutGitContext(&installCmdOptions)

	err = ensureIntegration(&installCmdOptions)
	if err != nil {
		eventSender.Fail(err.Error())
		return err, ""
	}

	// Need check if we want support not in cluster mode with Product owner
	installCmdOptions.Kube.InCluster = true

	questionnaire.AskAboutSyncOptions(&installCmdOptions)

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

	var kind, name string
	err, kind, name, manifest := templates.Install(&installOptions)

	if err != nil {
		msg := fmt.Sprintf("Argo agent installation resource \"%s\" with name \"%s\" finished with error , reason: %v ", kind, name, err)
		eventSender.Fail(err.Error())
		return errors.New(msg), ""
	}

	eventSender.Success("")

	logger.Success(fmt.Sprintf("Argo agent installation finished successfully to namespace \"%s\"", kubeOptions.Namespace))
	logger.Success(fmt.Sprintf("Gitops view: \"%s/gitops\"", installCmdOptions.Codefresh.Host))
	logger.Success(fmt.Sprintf("Documentation: \"%s\"", "https://codefresh.io/docs/docs/ci-cd-guides/gitops-deployments/"))
	return nil, manifest
}

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

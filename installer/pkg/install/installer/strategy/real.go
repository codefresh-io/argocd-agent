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
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/fatih/structs"
)

const (
	SUCCESS = "Success"
	FAILED  = "Failed"
)

type RealArgoCdAgentInstaller struct {
}

func (installer *RealArgoCdAgentInstaller) Install(installCmdOptions install.InstallCmdOptions) (error, string) {
	var err error
	// should be in beg for show correct events
	_ = questionnaire.AskAboutCodefreshCredentials(&installCmdOptions)

	holder.ApiHolder = codefresh.Api{
		Token:       installCmdOptions.Codefresh.Token,
		Host:        installCmdOptions.Codefresh.Host,
		Integration: installCmdOptions.Codefresh.Integration,
	}

	kubeConfigPath := installCmdOptions.Kube.ConfigPath
	kubeOptions := installCmdOptions.Kube
	if kubeOptions.Namespace == "" {
		kubeOptions.Namespace = "argocd"
	}

	if installCmdOptions.Agent.Interactive {
		_ = questionnaire.AskAboutKubeContext(&installCmdOptions)
	}

	kubeClient, err := kube.New(&kube.Options{
		ContextName:      kubeOptions.Context,
		Namespace:        kubeOptions.Namespace,
		PathToKubeConfig: kubeConfigPath,
		InCluster:        kubeOptions.InCluster,
	})

	if err != nil {
		return err, ""
	}

	if installCmdOptions.Agent.Interactive {
		_ = questionnaire.AskAboutNamespace(&installCmdOptions, kubeClient)
	}

	argoServerSvc, err := kubeClient.GetArgoServerSvc(kubeOptions.Namespace)
	if err != nil {
		msg := fmt.Sprintf("We didn't find ArgoCD on \"%s/%s\"", installCmdOptions.Kube.ClusterName, kubeOptions.Namespace)
		sendArgoAgentInstalledEvent(FAILED, msg)
		return errors.New(msg), ""
	}
	if kube.IsLoadBalancer(argoServerSvc) {
		balancerHost, _ := kubeClient.GetLoadBalancerHost(argoServerSvc)
		if balancerHost != "" {
			installCmdOptions.Argo.Host = balancerHost
		}
	}

	kubeOptions = installCmdOptions.Kube

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

	var kind, name string
	err, kind, name, manifest := templates.Install(&installOptions)

	if err != nil {
		msg := fmt.Sprintf("Argo agent installation resource \"%s\" with name \"%s\" finished with error , reason: %v ", kind, name, err)
		sendArgoAgentInstalledEvent(FAILED, msg)
		return errors.New(msg), ""
	}

	helper.ShowSummary(&installCmdOptions)

	sendArgoAgentInstalledEvent(SUCCESS, "")

	logger.Success(fmt.Sprintf("Argo agent installation finished successfully to namespace \"%s\"", kubeOptions.Namespace))
	logger.Success(fmt.Sprintf("Gitops view: \"%s/gitops\"", installCmdOptions.Codefresh.Host))
	logger.Success(fmt.Sprintf("Documentation: \"%s\"", "https://codefresh.io/docs/docs/ci-cd-guides/gitops-deployments/"))

	return nil, manifest
}

func sendArgoAgentInstalledEvent(status string, reason string) {
	props := make(map[string]string)
	props["status"] = status
	props["reason"] = reason
	_ = holder.ApiHolder.SendEvent("agent.installed", props)
}

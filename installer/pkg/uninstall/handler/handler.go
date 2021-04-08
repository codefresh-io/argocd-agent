package handler

import (
	"errors"
	"fmt"
	cfEventSender "github.com/codefresh-io/argocd-listener/installer/pkg/cfeventsender"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/questionnaire"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/codefresh-io/argocd-listener/installer/pkg/uninstall"
	"github.com/fatih/structs"
)

type UninstallHandler struct {
	cmdOptions  uninstall.CmdOptions
	eventSender *cfEventSender.CfEventSender
}

var uninstallHandler *UninstallHandler

func New(cmdOptions uninstall.CmdOptions) *UninstallHandler {
	if uninstallHandler == nil {
		eventSender := cfEventSender.New(cfEventSender.EVENT_AGENT_UNINSTALL)
		uninstallHandler = &UninstallHandler{cmdOptions, eventSender}
	}
	return uninstallHandler
}

func (uninstallHandler *UninstallHandler) Run() error {
	kubeOptions := uninstallHandler.cmdOptions.Kube

	_ = questionnaire.AskAboutKubeContext(&kubeOptions)

	kubeClient, err := kube.New(&kube.Options{
		ContextName:      kubeOptions.Context,
		Namespace:        kubeOptions.Namespace,
		PathToKubeConfig: uninstallHandler.cmdOptions.Kube.ConfigPath,
		InCluster:        kubeOptions.InCluster,
	})

	if err != nil {
		panic(err)
	}

	_ = questionnaire.AskAboutNamespace(&kubeOptions, kubeClient, false)

	uninstallOptions := templates.DeleteOptions{
		Templates:        kubernetes.TemplatesMap(),
		TemplateValues:   structs.Map(uninstallHandler.cmdOptions),
		Namespace:        kubeOptions.Namespace,
		KubeClientSet:    kubeClient.GetClientSet(),
		KubeCrdClientSet: kubeClient.GetCrdClientSet(),
	}

	var kind, name string
	err, kind, name = templates.Delete(&uninstallOptions)

	if err != nil {
		msg := fmt.Sprintf("Argo agent uninstallation resource \"%s\" with name \"%s\" finished with error , reason: %v ", kind, name, err)
		uninstallHandler.eventSender.Fail(msg)
		return errors.New(msg)
	}

	uninstallHandler.eventSender.Success("")

	logger.Success(fmt.Sprintf("Argo agent uninstallation finished successfully to namespace \"%s\"", kubeOptions.Namespace))
	return nil
}

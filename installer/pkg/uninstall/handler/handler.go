package handler

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/statuses"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates"
	"github.com/codefresh-io/argocd-listener/installer/pkg/templates/kubernetes"
	"github.com/codefresh-io/argocd-listener/installer/pkg/uninstall"
	"github.com/fatih/structs"
)

type UninstallHandler struct {
	cmdOptions uninstall.CmdOptions
}

var uninstallHandler *UninstallHandler

func New(cmdOptions uninstall.CmdOptions) *UninstallHandler {
	if uninstallHandler == nil {
		uninstallHandler = &UninstallHandler{cmdOptions}
	}
	return uninstallHandler
}

func (uninstallHandler *UninstallHandler) Run() error {
	kubeOptions := uninstallHandler.cmdOptions.Kube

	if uninstallHandler.cmdOptions.Kube.Context == "" {
		contexts, err := kube.GetAllContexts(uninstallHandler.cmdOptions.Kube.ConfigPath)
		if err != nil {
			return err
		}

		err, selectedContext := prompt.Select(contexts, "Select Kubernetes context")
		if err != nil {
			return err
		}
		kubeOptions.Context = selectedContext
	}

	kubeClient, err := kube.New(&kube.Options{
		ContextName:      kubeOptions.Context,
		Namespace:        kubeOptions.Namespace,
		PathToKubeConfig: uninstallHandler.cmdOptions.Kube.ConfigPath,
		InCluster:        kubeOptions.InCluster,
	})

	if err != nil {
		panic(err)
	}

	namespaces, err := kubeClient.GetNamespaces()
	if err != nil {
		err = prompt.InputWithDefault(&kubeOptions.Namespace, "Kubernetes namespace to uninstall", "default")
		if err != nil {
			return err
		}
	} else {
		err, selectedNamespace := prompt.Select(namespaces, "Select Kubernetes namespace to uninstall")
		if err != nil {
			return err
		}
		kubeOptions.Namespace = selectedNamespace
	}

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
		sendArgoAgentUninstalledEvent(statuses.FAILED, msg)
		return errors.New(msg)
	}

	sendArgoAgentUninstalledEvent(statuses.SUCCESS, "")

	logger.Success(fmt.Sprintf("Argo agent uninstallation finished successfully to namespace \"%s\"", kubeOptions.Namespace))
	return nil
}

func sendArgoAgentUninstalledEvent(status string, reason string) {
	props := make(map[string]string)
	props["status"] = status
	props["reason"] = reason
	_ = holder.ApiHolder.SendEvent("agent.uninstalled", props)
}

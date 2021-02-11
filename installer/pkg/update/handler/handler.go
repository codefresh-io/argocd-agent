package handler

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/kube"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/obj/kubeobj"
	"github.com/codefresh-io/argocd-listener/installer/pkg/questionnaire"
	"github.com/codefresh-io/argocd-listener/installer/pkg/update"
	"github.com/codefresh-io/argocd-listener/installer/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type UpdateHandler struct {
	cmdOptions update.CmdOptions
	version    string
}

var updateHandler *UpdateHandler

func New(cmdOptions update.CmdOptions, version string) *UpdateHandler {
	if updateHandler == nil {
		updateHandler = &UpdateHandler{cmdOptions, version}
	}
	return updateHandler
}

func (updateHandler *UpdateHandler) Run() error {
	var err error

	kubeConfigPath := updateHandler.cmdOptions.Kube.ConfigPath
	kubeOptions := updateHandler.cmdOptions.Kube

	_ = questionnaire.AskAboutKubeContext(&kubeOptions)

	kubeClient, err := kube.New(&kube.Options{
		ContextName:      kubeOptions.Context,
		Namespace:        kubeOptions.Namespace,
		PathToKubeConfig: kubeConfigPath,
		InCluster:        kubeOptions.InCluster,
	})
	if err != nil {
		return err
	}

	_ = questionnaire.AskAboutNamespace(&kubeOptions, kubeClient, false)
	err = updateDeploymentWithNewVersion(kubeClient.GetClientSet(), kubeOptions.Namespace, updateHandler.cmdOptions.Codefresh.Suffix, updateHandler.version)

	if err != nil {
		return errors.New(fmt.Sprintf("Argo agent update finished with error , reason: %v ", err))
	}

	logger.Success(fmt.Sprintf("Argo agent update finished successfully to namespace \"%s\"", kubeOptions.Namespace))

	return nil
}

func updateDeploymentWithNewVersion(clientSet *kubernetes.Clientset, namespace string, suffix string, version string) error {
	deploymentList, err := kubeobj.GetDeployments(clientSet, namespace, "app=cf-argocd-agent"+suffix)

	if err != nil {
		return errors.New(fmt.Sprintf("Argo agent update finished with error , reason: %v ", err))
	}

	if len(deploymentList.Items) == 0 {
		return errors.New("Argo agent failed to update because no deployments were found")
	}

	deployment := &deploymentList.Items[0]

	envs := deployment.Spec.Template.Spec.Containers[0].Env

	newEnvs := make([]v1.EnvVar, 0)

	for _, env := range envs {
		if env.Name == "AGENT_VERSION" {
			env.Value = util.ResolvePackageVersion(version)
		}
		newEnvs = append(newEnvs, env)
	}

	deployment.Spec.Template.Spec.Containers[0].Env = newEnvs

	_, err = kubeobj.UpdateDeployment(clientSet, deployment, namespace)

	return err
}

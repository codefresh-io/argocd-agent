package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
)

func Run(argoOptions *install.ArgoOptions) error {
	logger.Info("\nTesting requirements")
	logger.Info("--------------------")
	defer logger.Info("--------------------\n")

	credentialsMsg := "checking argocd credentials..."
	projectsMsg := "checking argocd projects accessibility..."
	applicationsMsg := "checking argocd applications accessibility..."
	var err error

	err = checkArgoCredentials(argoOptions)
	if err != nil {
		logger.FailureTest(credentialsMsg)
		return err
	}
	logger.SuccessTest(credentialsMsg)

	err = checkProjects()
	if err != nil {
		logger.FailureTest(projectsMsg)
		return err
	}
	logger.SuccessTest(projectsMsg)

	err = checkApplications()
	if err != nil {
		logger.FailureTest(applicationsMsg)
		return err
	}
	logger.SuccessTest(applicationsMsg)
	return nil
}

func checkArgoCredentials(argoOptions *install.ArgoOptions) error {
	var err error
	token := argoOptions.Token
	if token == "" {
		token, err = argo.GetToken(argoOptions.Username, argoOptions.Password, argoOptions.Host)
		if err == nil {
			store.SetArgo(token, argoOptions.Host)
		}
	} else {
		store.SetArgo(token, argoOptions.Host)
		err = argo.GetInstance().CheckToken()
	}
	return err
}

func checkProjects() error {
	_, err := argo.GetProjectsWithCredentialsFromStorage()
	return err
}

func checkApplications() error {
	_, err := argo.GetInstance().GetApplicationsWithCredentialsFromStorage()
	return err
}

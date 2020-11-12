package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
)

func Run(argoOptions *install.ArgoOptions) error {
	info("\nTesting requirements")
	info("--------------------")
	defer info("--------------------\n")
	
	credentialsMsg := "checking argocd credentials..."
	projectsMsg := "checking argocd projects accessibility..."
	applicationsMsg := "checking argocd applications accessibility..."
	var err error
	var token string

	token, err = checkArgoCredentials(argoOptions)
	if err != nil {
		failure(credentialsMsg)
		return err
	}
	success(credentialsMsg)
	store.SetArgo(token, argoOptions.Host)

	err = checkProjects()
	if err != nil {
		failure(projectsMsg)
		return err
	}
	success(projectsMsg)

	err = checkApplications()
	if err != nil {
		failure(applicationsMsg)
		return err
	}
	success(applicationsMsg)
	return nil
}

func checkArgoCredentials(argoOptions *install.ArgoOptions) (string, error){
	var err error
	token := argoOptions.Token
	if token == "" {
		token, err = argo.GetToken(argoOptions.Username, argoOptions.Password, argoOptions.Host)
	}else{
		err = argo.CheckToken(token, argoOptions.Host)
	}
	return token, err
}

func checkProjects() error {
	_, err := argo.GetProjectsWithCredentialsFromStorage()
	return err
}

func checkApplications() error{
	_, err := argo.GetInstance().GetApplicationsWithCredentialsFromStorage()
	return err
}

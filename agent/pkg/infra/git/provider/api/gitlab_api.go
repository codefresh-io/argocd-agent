package api

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"github.com/xanzy/go-gitlab"
)

type (
	gitlabApi struct {
		git *gitlab.Client
	}

	GitlabApi interface {
		ListProjects(page int) (error, []*gitlab.Project)
		RetrieveAvatar(email string) (error, string)
		GetCommit(projectId int, revision string) (error, *gitlab.Commit)
	}
)

func NewGitlabApi() GitlabApi {
	context := &store.GetStore().Git.Context

	var clientOptions gitlab.ClientOptionFunc
	if context.Spec.Data.Auth.ApiHost != "" {
		logger.GetLogger().Infof("Override gitlab uri with custom one %s", context.Spec.Data.Auth.ApiHost)
		clientOptions = gitlab.WithBaseURL(context.Spec.Data.Auth.ApiHost)
	} else {
		//TODO: verify why Haim has empty api host, if it is not
		logger.GetLogger().Infof("Override gitlab uri with custom one, hk  %s", "https://gitlab.hosts-app.com/api/v4/")
		clientOptions = gitlab.WithBaseURL("https://gitlab.hosts-app.com/api/v4/")
	}

	git, err := gitlab.NewOAuthClient(context.Spec.Data.Auth.Password, clientOptions)
	if err != nil {
		logger.GetLogger().Errorf("Cant initialize gitlab oauth client %s because %v ", context.Spec.Data.Auth.ApiHost, err.Error())
	}

	maskedPassword := util.MaskLeft(context.Spec.Data.Auth.Password)
	logger.GetLogger().Infof("Initializing gitlab client, host: %s, password %s", git.BaseURL().String(), maskedPassword)

	return &gitlabApi{git}
}

// TODO: contibute this api to go-gitlab library
func (gitlabApi *gitlabApi) RetrieveAvatar(email string) (error, string) {
	opts := struct {
		Email string `url:"email,omitempty"`
	}{Email: email}
	req, err := gitlabApi.git.NewRequest("GET", "/avatar", opts, nil)
	if err != nil {
		return err, ""
	}

	avatar := struct {
		Url string `json:"avatar_url"`
	}{}

	_, err = gitlabApi.git.Do(req, &avatar)
	if err != nil {
		return err, ""
	}

	return nil, avatar.Url
}

func (gitlabApi *gitlabApi) ListProjects(page int) (error, []*gitlab.Project) {
	owner := true
	listProjectOptions := &gitlab.ListProjectsOptions{
		Membership:  &owner,
		ListOptions: gitlab.ListOptions{Page: page},
	}
	projects, _, err := gitlabApi.git.Projects.ListProjects(listProjectOptions)
	return err, projects
}

func (gitlabApi *gitlabApi) GetCommit(projectId int, revision string) (error, *gitlab.Commit) {
	commit, _, err := gitlabApi.git.Commits.GetCommit(projectId, revision)
	return err, commit
}

package api

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/xanzy/go-gitlab"
)

type (
	gitlabApi struct {
		git *gitlab.Client
	}

	GitlabApi interface {
		ListProjects(search string) (error, []*gitlab.Project)
		RetrieveAvatar(email string) (error, string)
		GetCommit(projectId int, revision string) (error, *gitlab.Commit)
	}
)

func NewGitlabApi() GitlabApi {
	context := &store.GetStore().Git.Context
	git, _ := gitlab.NewOAuthClient(context.Spec.Data.Auth.Password, gitlab.WithBaseURL(context.Spec.Data.Auth.ApiHost))
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

func (gitlabApi *gitlabApi) ListProjects(search string) (error, []*gitlab.Project) {
	owner := true
	listProjectOptions := &gitlab.ListProjectsOptions{
		Search: &search,
		Owned:  &owner,
	}
	projects, _, err := gitlabApi.git.Projects.ListProjects(listProjectOptions)
	return err, projects
}

func (gitlabApi *gitlabApi) GetCommit(projectId int, revision string) (error, *gitlab.Commit) {
	commit, _, err := gitlabApi.git.Commits.GetCommit(projectId, revision)
	return err, commit
}

package provider

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/thoas/go-funk"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type (
	Gitlab struct {
		git *gitlab.Client
	}
)

var gitlabInstance *Gitlab

func NewGitlabProvider() GitProvider {
	if gitlabInstance == nil {
		context := &store.GetStore().Git.Context
		git, _ := gitlab.NewOAuthClient(context.Spec.Data.Auth.Password, gitlab.WithBaseURL(context.Spec.Data.Auth.ApiHost))
		gitlabInstance = &Gitlab{git: git}
	}
	return gitlabInstance
}

// TODO: contibute this api to go-gitlab library
func (gitlabInstance *Gitlab) retrieveAvatar(email string) (error, string) {
	opts := struct {
		Email string `url:"email,omitempty"`
	}{Email: email}
	req, err := gitlabInstance.git.NewRequest("GET", "/avatar", opts, nil)
	if err != nil {
		return err, ""
	}

	avatar := struct {
		Url string `json:"avatar_url"`
	}{}

	_, err = gitlabInstance.git.Do(req, &avatar)
	if err != nil {
		return err, ""
	}

	return nil, avatar.Url
}

func (gitlabInstance *Gitlab) GetCommitByRevision(repoUrl string, revision string) (error, *service.ResourceCommit) {
	parts := strings.Split(repoUrl, "/")
	// p.kostohrys/test
	if len(parts) != 2 {
		return errors.New("wrong amount of arguments"), nil
	}
	owner := true
	listProjectOptions := &gitlab.ListProjectsOptions{
		Search: &parts[1],
		Owned:  &owner,
	}

	projects, _, err := gitlabInstance.git.Projects.ListProjects(listProjectOptions)

	if err != nil {
		return err, nil
	}

	proj, ok := funk.Find(projects, func(proj *gitlab.Project) bool {
		return proj.PathWithNamespace == repoUrl
	}).(*gitlab.Project)

	if !ok {
		return errors.New("failed to find gitlab project"), nil
	}

	commit, _, err := gitlabInstance.git.Commits.GetCommit(proj.ID, revision)

	if err != nil {
		return err, nil
	}

	err, avatar := gitlabInstance.retrieveAvatar(commit.AuthorEmail)

	if err != nil {
		return err, nil
	}

	result := &service.ResourceCommit{
		Message: &commit.Message,
		Sha:     &revision,
		Avatar:  &avatar,
	}

	result.Link = &commit.WebURL

	return nil, result
}

func (gitlab *Gitlab) GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops) {
	defaultGitInfo := codefreshSdk.Gitops{
		Comitters: []codefreshSdk.User{},
		Prs:       []codefreshSdk.Annotation{},
		Issues:    []codefreshSdk.Annotation{},
	}
	return nil, &defaultGitInfo
}

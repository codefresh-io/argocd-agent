package provider

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/git/provider/api"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/thoas/go-funk"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type (
	Gitlab struct {
		api api.GitlabApi
	}
)

var gitlabInstance *Gitlab

func NewGitlabProvider() GitProvider {
	if gitlabInstance == nil {
		gitlabInstance = &Gitlab{api: api.NewGitlabApi()}
	}
	return gitlabInstance
}

func (gitlabInstance *Gitlab) GetCommitByRevision(repoUrl string, revision string) (error, *service.ResourceCommit) {
	logger.GetLogger().Infof("Start handle get commit by revision for repo %s and revision %s", repoUrl, revision)
	parts := strings.Split(repoUrl, "/")
	// p.kostohrys/test
	if len(parts) != 2 {
		return errors.New("wrong amount of arguments"), nil
	}

	err, projects := gitlabInstance.api.ListProjects(parts[1])

	if err != nil {
		return err, nil
	}

	proj, ok := funk.Find(projects, func(proj *gitlab.Project) bool {
		return proj.PathWithNamespace == repoUrl
	}).(*gitlab.Project)

	if !ok {
		return errors.New("failed to find gitlab project"), nil
	}

	err, commit := gitlabInstance.api.GetCommit(proj.ID, revision)

	if err != nil {
		return err, nil
	}

	err, avatar := gitlabInstance.api.RetrieveAvatar(commit.AuthorEmail)

	if err != nil {
		avatar = ""
		logger.GetLogger().Infof("Setup empty avatar for user %s because of error", commit.AuthorEmail)
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

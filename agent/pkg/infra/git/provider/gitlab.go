package provider

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/git/provider/api"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/thoas/go-funk"
	"github.com/xanzy/go-gitlab"
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

	var page = 0

	for {
		err, projects := gitlabInstance.api.ListProjects(page)

		if err != nil {
			return err, nil
		}

		if len(projects) == 0 {
			break
		}

		page++

		proj, ok := funk.Find(projects, func(proj *gitlab.Project) bool {
			return proj.HTTPURLToRepo == repoUrl || proj.SSHURLToRepo == repoUrl
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

	return errors.New(fmt.Sprintf("Project with name %s not found in gitlab", repoUrl)), nil
}

func (gitlab *Gitlab) GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops) {
	defaultGitInfo := codefreshSdk.Gitops{
		Comitters: []codefreshSdk.User{},
		Prs:       []codefreshSdk.Annotation{},
		Issues:    []codefreshSdk.Annotation{},
	}
	return nil, &defaultGitInfo
}

package provider

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/git/provider/api"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
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

func (gitlabInstance *Gitlab) getProject(repoUrl string) (error, *gitlab.Project) {
	var page = 1

	for {
		err, projects := gitlabInstance.api.ListProjects(page)

		if err != nil {
			return err, nil
		}

		if len(projects) == 0 {
			return errors.New("failed to find gitlab project"), nil
		}

		page++

		foundedProject := funk.Find(projects, func(proj *gitlab.Project) bool {
			return proj.HTTPURLToRepo == repoUrl || proj.SSHURLToRepo == repoUrl
		})

		if foundedProject != nil {
			logger.GetLogger().Infof("found gitlab project to map %v", foundedProject)
		}

		proj, ok := foundedProject.(*gitlab.Project)

		if !ok {
			continue
		}

		return nil, proj
	}
}

func (gitlabInstance *Gitlab) GetCommitByRevision(repoUrl string, revision string) (error, *ResourceCommit) {
	logger.GetLogger().Infof("Start handle get commit by revision for repo %s and revision %s", repoUrl, revision)

	err, project := gitlabInstance.getProject(repoUrl)
	if err != nil {
		return err, nil
	}

	err, commit := gitlabInstance.api.GetCommit(project.ID, revision)

	if err != nil {
		return err, nil
	}

	err, avatar := gitlabInstance.api.RetrieveAvatar(commit.AuthorEmail)

	if err != nil {
		avatar = ""
		logger.GetLogger().Infof("Setup empty avatar for user %s because of error", commit.AuthorEmail)
	}

	result := &ResourceCommit{
		Message: &commit.Message,
		Sha:     &revision,
		Avatar:  &avatar,
	}

	result.Link = &commit.WebURL

	return nil, result
}

func (gitlabInstance *Gitlab) GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops) {
	logger.GetLogger().Infof("Start handle get manifest  for repo %s and revision %s", repoUrl, revision)

	defaultGitInfo := codefreshSdk.Gitops{
		Comitters: []codefreshSdk.GitopsUser{},
		Prs:       []codefreshSdk.Annotation{},
		Issues:    []codefreshSdk.Annotation{},
	}

	err, project := gitlabInstance.getProject(repoUrl)
	if err != nil {
		return err, &defaultGitInfo
	}

	err, commit := gitlabInstance.api.GetCommit(project.ID, revision)
	if err != nil {
		return err, &defaultGitInfo
	}

	commits := []*gitlab.Commit{
		commit,
	}

	err, committers := gitlabInstance.api.GetComittersByCommits(commits)
	if err != nil {
		return err, &defaultGitInfo
	}

	err, prs := gitlabInstance.api.GetPrsByCommits(project.ID, commits)
	if err != nil {
		return err, &defaultGitInfo
	}

	gitInfo := codefreshSdk.Gitops{
		Comitters: committers,
		Prs:       prs,
		Issues:    []codefreshSdk.Annotation{},
	}

	return nil, &gitInfo
}

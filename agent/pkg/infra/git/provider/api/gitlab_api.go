package api

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
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
		GetCommitsBySha(projectId int, revision string) (error, []*gitlab.Commit)
		GetComittersByCommits(commits []*gitlab.Commit) (error, []codefreshSdk.GitopsUser)
		GetPrsByCommits(projectId int, commits []*gitlab.Commit) (error, []codefreshSdk.Annotation)
	}
)

func NewGitlabApi() GitlabApi {
	context := &store.GetStore().Git.Context

	var clientOptions gitlab.ClientOptionFunc
	if context.Spec.Data.Auth.ApiURL != "" {
		logger.GetLogger().Infof("Override gitlab uri with custom one %s", context.Spec.Data.Auth.ApiURL)
		clientOptions = gitlab.WithBaseURL(context.Spec.Data.Auth.ApiURL)
	}

	git, err := gitlab.NewOAuthClient(context.Spec.Data.Auth.Password, clientOptions)
	if err != nil {
		logger.GetLogger().Errorf("Cant initialize gitlab oauth client %s because %v ", context.Spec.Data.Auth.ApiURL, err.Error())
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

func (gitlabApi *gitlabApi) GetCommitsBySha(projectId int, revision string) (error, []*gitlab.Commit) {
	commitsOption := &gitlab.ListCommitsOptions{
		RefName: &revision,
	}
	commit, _, err := gitlabApi.git.Commits.ListCommits(projectId, commitsOption)
	return err, commit
}

func (gitlabApi *gitlabApi) GetComittersByCommits(commits []*gitlab.Commit) (error, []codefreshSdk.GitopsUser) {
	committers := make(map[string]*codefreshSdk.GitopsUser)

	for _, commit := range commits {
		if committers[commit.CommitterName] == nil {
			_, avatar := gitlabApi.RetrieveAvatar(commit.CommitterEmail)
			committers[commit.CommitterName] = &codefreshSdk.GitopsUser{
				Name:   commit.CommitterName,
				Avatar: avatar,
			}
		}
	}

	users := []codefreshSdk.GitopsUser{}

	for _, value := range committers {
		users = append(users, *value)
	}

	return nil, users
}

func (gitlabApi *gitlabApi) GetPrsByCommits(projectId int, commits []*gitlab.Commit) (error, []codefreshSdk.Annotation) {
	opt := &gitlab.ListProjectMergeRequestsOptions{}
	prs, _, err := gitlabApi.git.MergeRequests.ListProjectMergeRequests(projectId, opt)
	if err != nil {
		return err, nil
	}
	pullRequests := []codefreshSdk.Annotation{}

	for _, pr := range prs {
		mergeCommitSHA := pr.MergeCommitSHA
		if mergeCommitSHA == "" {
			continue
		}
		for _, commit := range commits {
			if commit.ID == "" {
				continue
			}
			if commit.ID == mergeCommitSHA {

				pullRequests = append(pullRequests, codefreshSdk.Annotation{
					Key:   pr.Title,
					Value: pr.WebURL,
				})
			}
		}
	}
	return nil, pullRequests
}

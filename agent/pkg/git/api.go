package git

import (
	"context"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Api struct {
	Token  string
	Client *github.Client //github.Client
	Owner  string
	Repo   string
	Ctx    context.Context
}

var api *Api

func GetInstance() *Api {
	if api != nil {
		return api
	}
	gitConfig := store.GetStore().Git
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitConfig.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	api = &Api{
		Token:  gitConfig.Token,
		Ctx:    ctx,
		Client: client,
		Owner:  "olegz-codefresh", //todo - get this from config
		Repo:   "argo",            //todo - get this from config
	}
	return api
}

func (a *Api) GetCommitsBySha(sha string) (error, []*github.RepositoryCommit) {
	revisionCommit, _, err := api.Client.Repositories.GetCommit(api.Ctx, api.Owner, api.Repo, sha)
	if err != nil {
		return err, nil
	}
	commits := []*github.RepositoryCommit{revisionCommit}
	if len(revisionCommit.Parents) > 0 {
		for i := 0; i < len(revisionCommit.Parents); i++ {
			commitInfo, _, err := api.Client.Repositories.GetCommit(api.Ctx, api.Owner, api.Repo, *revisionCommit.Parents[i].SHA)
			commits = append(commits, commitInfo)
			if err != nil {
				return err, nil
			}
		}
	}

	return nil, commits
}

func (a *Api) GetCommittersByCommits(commits []*github.RepositoryCommit) (error, []User) {
	committers := []User{}
	committersSet := make(map[string]bool)
	for _, commit := range commits {
		author := commit.Author
		if author != nil {
			_, exists := committersSet[*author.Login]
			if exists != true {
				committersSet[*author.Login] = true
				committers = append(committers, User{
					Name:   author.String(),
					Avatar: "",
				})
			}
		}
	}

	return nil, committers
}

func (a *Api) GetPullRequestsByCommits(commits []*github.RepositoryCommit) (error, []*github.PullRequest) {
	allPullRequests, _, err := api.Client.PullRequests.List(api.Ctx, api.Owner, api.Repo, &github.PullRequestListOptions{State: "all"})
	if err != nil {
		return err, nil
	}

	pullRequests := []*github.PullRequest{}
	for _, pr := range allPullRequests {
		mergeCommitSHA := pr.MergeCommitSHA
		if mergeCommitSHA == nil {
			continue
		}
		for _, commit := range commits {
			if commit.SHA == nil {
				continue
			}
			if *commit.SHA == *mergeCommitSHA {
				pullRequests = append(pullRequests, pr)
			}
		}
	}
	return nil, pullRequests
}

func (a *Api) GetIssuesByPRs(pullRequests []*github.PullRequest) (error, []*github.Issue) {
	allIssues := []*github.Issue{}

	for _, prs := range pullRequests {
		issues, _, err := api.Client.Issues.Get(api.Ctx, api.Owner, api.Repo, *prs.Number)
		if err != nil {
			return err, nil
		}
		allIssues = append(allIssues, issues)
	}

	return nil, allIssues
}

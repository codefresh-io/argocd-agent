package git

import (
	"context"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Api struct {
	Token  string
	Client *github.Client //github.Client
	Owner  string
	Repo   string
	Ctx context.Context
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
		Ctx: ctx,
		Client: client,
		Owner:  "olegz-codefresh",
		Repo:   "argo",
	}
	return api
}

func (a *Api) GetCommitsBySha(sha string) (error, []*github.RepositoryCommit) {
	// @todo - wtf with pointers

	revisionCommit, _, err := api.Client.Repositories.GetCommit(api.Ctx, api.Owner, api.Repo, sha)
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

	if err != nil {
		return err, nil
	}
	return nil, commits
}

func (a *Api) GetCommittersByCommits(commits []*github.RepositoryCommit) (error, []*github.User) {
	// @todo - wtf with pointers
	committers := []*github.User{}
	for i := 0; i < len(commits); i++ {
		author := commits[i].Author
		committers = append(committers, author)
	}

	return nil, committers
}

func (a *Api) GetPullRequestsByCommits(commits []*github.RepositoryCommit) (error, *github.PullRequest) {
	// @todo - wtf with pointers
	allPullRequests, _, err := api.Client.PullRequests.List(api.Ctx, api.Owner, api.Repo, &github.PullRequestListOptions{State: "all"})

	pullRequests := []*github.RepositoryCommit{}
	if err != nil {
		return err, nil
	}
	fmt.Println(pullRequests)

	//allSha := []string
	//for i := 0; i < len(commits); i++ {
	//	allSha := append(allSha, *commits[i].SHA)
	//}
	//
	//for i := 0; i < len(allPullRequests); i++ {
	//	mergeCommitSHA := allPullRequests[i].MergeCommitSHA
	//	fmt.Println(MergeCommitSHA)
	//}
	//for i := 0; i < len(*commits); i++ {
	//	author := (*commits)[i].Author
	//	committers = append(committers, author)
	//}

	return nil, allPullRequests[0]
}

func (a *Api) GetIssuesByPRs(pullRequest *github.PullRequest) (error, interface{}) {
	// @todo - remove this sh*t
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: api.Token}))
	Client := github.NewClient(tc)

	return nil, Client
}
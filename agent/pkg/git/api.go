package git

import (
	"context"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	_ "golang.org/x/oauth2"
)

type Api struct {
	Token  string
	Client interface{} //github.Client
	Owner  string
	Repo   string
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
		Client: client,
		Owner:  "olegz-codefresh",
		Repo:   "argo",
	}
	return api
}

func (a *Api) GetCommitsBySha(sha string) (error, []github.RepositoryCommit) {
	// @todo - remove this sh*t
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: api.Token}))
	Client := github.NewClient(tc)

	revisionCommit, _, err := Client.Repositories.GetCommit(ctx, api.Owner, api.Repo, sha)
	commits := []github.RepositoryCommit{*revisionCommit}
	if len(revisionCommit.Parents) > 0 {
		for i := 0; i < len(revisionCommit.Parents); i++ {
			commitInfo, _, err := Client.Repositories.GetCommit(ctx, api.Owner, api.Repo, *revisionCommit.Parents[i].SHA)
			commits = append(commits, *commitInfo)
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

func (a *Api) GetCommittersByCommits(commits *[]github.RepositoryCommit) (error, []*github.User) {
	// @todo - wtf with pointers
	committers := []*github.User{}
	for i := 0; i < len(*commits); i++ {
		author := (*commits)[i].Author
		committers = append(committers, author)
	}

	return nil, committers
}

func (a *Api) GetPullRequestsByCommits(commits *[]github.RepositoryCommit) (error, interface{}) {
	// @todo - remove this sh*t
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: api.Token}))
	Client := github.NewClient(tc)

	//for i := 0; i < len(*commits); i++ {
	//	author := (*commits)[i].Author
	//	committers = append(committers, author)
	//}

	return nil, Client
}

func (a *Api) GetIssuesByPRs(sha string) (error, interface{}) {
	// @todo - remove this sh*t
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: api.Token}))
	Client := github.NewClient(tc)

	return nil, Client
}

package transform

import (
	"encoding/json"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/git"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/mitchellh/mapstructure"
	"sort"
)

func initDeploymentsStatuses(applicationName string) map[string]string {
	statuses := make(map[string]string)
	resourceTree, _ := argo.GetResourceTree(applicationName)
	for _, node := range resourceTree.Nodes {
		if node.Kind == "Deployment" || node.Kind == "Rollout" {
			if node.Health.Status == "" {
				statuses[node.Uid] = "Missing"
			} else {
				statuses[node.Uid] = node.Health.Status
			}

		}
	}
	return statuses
}

func prepareEnvironmentActivity(applicationName string) ([]codefresh2.EnvironmentActivity, error) {

	resource, err := argo.GetManagedResources(applicationName)
	if err != nil {
		return nil, err
	}

	statuses := initDeploymentsStatuses(applicationName)

	var activities []codefresh2.EnvironmentActivity

	for _, item := range resource.Items {
		if item.Kind == "Deployment" || item.Kind == "Rollout" {

			var targetState argo.ManagedResourceState
			err = json.Unmarshal([]byte(item.TargetState), &targetState)
			if err != nil {
				logger.GetLogger().Errorf("Failed to unmarshal \"TargetState\" to ManagedResourceState, reason %v", err)
				continue
			}

			var targetImages []string
			for _, container := range targetState.Spec.Template.Spec.Containers {
				targetImages = append(targetImages, container.Image)
			}

			var liveState argo.ManagedResourceState
			err = json.Unmarshal([]byte(item.LiveState), &liveState)
			if err != nil {
				logger.GetLogger().Errorf("Failed to unmarshal \"LiveState\" to ManagedResourceState, reason %v", err)
				continue
			}

			var liveImages []string
			for _, container := range liveState.Spec.Template.Spec.Containers {
				liveImages = append(liveImages, container.Image)
			}
			status := statuses[liveState.Metadata.Uid]
			activities = append(activities, codefresh2.EnvironmentActivity{
				Name:         item.Name,
				TargetImages: targetImages,
				Status:       status,
				LiveImages:   liveImages,
			})
		}
	}

	return activities, nil
}

func PrepareEnvironment(envItem map[string]interface{}) (error, *codefresh2.Environment) {

	var app argo.ArgoApplication
	err := mapstructure.Decode(envItem, &app)
	if err != nil {
		return err, nil
	}

	name := app.Metadata.Name
	historyList := app.Status.History
	revision := app.Status.OperationState.SyncResult.Revision
	repoUrl := app.Spec.Source.RepoURL

	resources, err := argo.GetResourceTreeAll(name)
	if err != nil {
		return err, nil
	}

	// we still need send env , even if we have problem with retrieve gitops info
	err, gitops := getGitoptsInfo(repoUrl, revision)

	if err != nil {
		logger.GetLogger().Errorf("Failed to retrieve manifest repo git information , reason: %v", err)
	}

	err, historyId := resolveHistoryId(historyList, app.Status.OperationState.SyncResult.Revision, name)

	if err != nil {
		return err, nil
	}

	activities, err := prepareEnvironmentActivity(name)
	if err != nil {
		return err, nil
	}

	env := codefresh2.Environment{
		HealthStatus: app.Status.Health.Status,
		SyncStatus:   app.Status.Sync.Status,
		SyncRevision: revision,
		Gitops:       *gitops,
		HistoryId:    historyId,
		Name:         name,
		Activities:   activities,
		Resources:    resources,
		RepoUrl:      repoUrl,
		FinishedAt:   app.Status.OperationState.FinishedAt,
	}

	return nil, &env

}

func resolveHistoryId(historyList []argo.ArgoApplicationHistoryItem, revision string, name string) (error, int64) {
	if historyList == nil {
		logger.GetLogger().Errorf("can`t find history id for application %s, because history list is empty", name)
		return nil, -1
	}

	sort.Slice(historyList, func(i, j int) bool {
		return historyList[i].Id > historyList[j].Id
	})

	for _, item := range historyList {
		if item.Revision == revision {
			return nil, item.Id
		}
	}
	return fmt.Errorf("can`t find history id for application %s", name), 0
}

func getGitoptsInfo(repoUrl string, revision string) (error, *git.Gitops) {
	defaultGitInfo := git.Gitops{
		Comitters: []git.User{},
		Prs:       []git.Annotation{},
		Issues:    []git.Annotation{},
	}
	err, gitClient := git.GetInstance(repoUrl)
	if err != nil {
		return err, &defaultGitInfo
	}

	err, commits := gitClient.GetCommitsBySha(revision)
	if err != nil {
		return err, &defaultGitInfo
	}

	err, comitters := gitClient.GetComittersByCommits(commits)
	if err != nil {
		return err, &defaultGitInfo
	}

	err, _, prs := gitClient.GetIssuesAndPrsByCommits(commits)
	if err != nil {
		return err, &defaultGitInfo
	}

	gitInfo := git.Gitops{
		Comitters: comitters,
		Prs:       prs,
		Issues:    []git.Annotation{},
	}

	return nil, &gitInfo
}

package transform

import (
	"encoding/json"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	git "github.com/codefresh-io/argocd-listener/agent/pkg/git"
	"github.com/mitchellh/mapstructure"
	"log"
	"sort"
)

type ArgoApplicationHistoryItem struct {
	Id       int64
	Revision string
}

type ArgoApplication struct {
	Status struct {
		Health struct {
			Status string
		}
		Sync struct {
			Status   string
			Revision string
		}
		History        []ArgoApplicationHistoryItem
		OperationState struct {
			FinishedAt string
			SyncResult struct {
				Revision string
			}
		}
	}
	Spec struct {
		Source struct {
			RepoURL string
		}
	}
	Metadata struct {
		Name string
	}
}

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

func prepareEnvironmentActivity(applicationName string) []codefresh2.EnvironmentActivity {

	resource := argo.GetManagedResources(applicationName)

	statuses := initDeploymentsStatuses(applicationName)

	var activities []codefresh2.EnvironmentActivity

	for _, item := range resource.Items {
		if item.Kind == "Deployment" || item.Kind == "Rollout" {

			var targetState argo.ManagedResourceState
			err := json.Unmarshal([]byte(item.TargetState), &targetState)
			if err != nil {
				log.Println(err.Error())
			}

			var targetImages []string
			for _, container := range targetState.Spec.Template.Spec.Containers {
				targetImages = append(targetImages, container.Image)
			}

			var liveState argo.ManagedResourceState
			err = json.Unmarshal([]byte(item.LiveState), &liveState)
			if err != nil {
				log.Println(err.Error())
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

	return activities
}

func PrepareEnvironment(envItem map[string]interface{}) (error, *codefresh2.Environment) {

	var app ArgoApplication
	err := mapstructure.Decode(envItem, &app)

	name := app.Metadata.Name
	revision := app.Status.OperationState.SyncResult.Revision
	if name != "app-olegz" {
		_env := codefresh2.Environment{}
		return nil, &_env
	}

	historyList := app.Status.History
	err, gitInfo := getGitObject(revision)
	fmt.Println(gitInfo)
	resources, err := argo.GetResourceTreeAll(name)
	// TODO: improve error handling
	if err != nil {
		return err, nil
	}

	err, historyId := resolveHistoryId(historyList, app.Status.OperationState.SyncResult.Revision)

	if err != nil {
		return err, nil
	}

	env := codefresh2.Environment{
		HealthStatus: app.Status.Health.Status,
		SyncStatus:   app.Status.Sync.Status,
		SyncRevision: revision,
		HistoryId:    historyId,
		Name:         name,
		Activities:   prepareEnvironmentActivity(name),
		Resources:    resources,
		RepoUrl:      app.Spec.Source.RepoURL,
		FinishedAt:   app.Status.OperationState.FinishedAt,
	}

	return nil, &env

}

func resolveHistoryId(historyList []ArgoApplicationHistoryItem, revision string) (error, int64) {
	sort.Slice(historyList, func(i, j int) bool {
		return historyList[i].Id > historyList[j].Id
	})

	for _, item := range historyList {
		if item.Revision == revision {
			return nil, item.Id
		}
	}
	return fmt.Errorf("can`t find history id"), 0
}

func getGitObject(revision string) (error, *codefresh2.GitInfo) {

	gitClient := git.GetInstance()

	err, commits := gitClient.GetCommitsBySha(revision)
	if err != nil { // @todo - maybe we have better idea
		return err, nil
	}

	err, committers := gitClient.GetCommittersByCommits(&commits)
	if err != nil { // @todo - maybe we have better idea
		return err, nil
	}

	err, prs := gitClient.GetCommittersByCommits(&commits)
	if err != nil { // @todo - maybe we have better idea
		return err, nil
	}

	fmt.Println(commits)
	fmt.Println(committers)
	fmt.Println(prs)
	gitObject := codefresh2.GitInfo{}
	return nil, &gitObject
}

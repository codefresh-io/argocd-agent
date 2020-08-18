package transform

import (
	"encoding/json"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/mitchellh/mapstructure"
	"log"
)

type ArgoApplication struct {
	Status struct {
		Health struct {
			Status string
		}
		Sync struct {
			Status   string
			Revision string
		}
		History []struct {
			Id int64
		}
		OperationState struct {
			FinishedAt string
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

func PrepareEnvironment(envItem map[string]interface{}) codefresh2.Environment {

	var app ArgoApplication
	err := mapstructure.Decode(envItem, &app)

	name := app.Metadata.Name
	historyList := app.Status.History

	resources, err := argo.GetResourceTreeAll(name)
	// TODO: improve error handling
	if err != nil {
		println(err)
	}

	env := codefresh2.Environment{
		HealthStatus: app.Status.Health.Status,
		SyncStatus:   app.Status.Sync.Status,
		SyncRevision: app.Status.Sync.Revision,
		HistoryId:    resolveHistoryId(historyList),
		Name:         name,
		Activities:   prepareEnvironmentActivity(name),
		Resources:    resources,
		RepoUrl:      app.Spec.Source.RepoURL,
		FinishedAt:   app.Status.OperationState.FinishedAt,
	}

	return env

}

func resolveHistoryId(historyList []struct {
	Id int64
}) int64 {
	if len(historyList) == 0 {
		return -1
	}
	return historyList[len(historyList)-1].Id
}

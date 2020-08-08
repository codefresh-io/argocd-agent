package argo

import (
	"encoding/json"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"log"
)

func initDeploymentsStatuses(applicationName string) map[string]string {
	statuses := make(map[string]string)
	resourceTree, _ := GetResourceTree(applicationName)
	for _, node := range resourceTree.Nodes {
		if node.Kind == "Deployment" {
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

	resource := GetManagedResources(applicationName)

	statuses := initDeploymentsStatuses(applicationName)

	var activities []codefresh2.EnvironmentActivity

	for _, item := range resource.Items {
		if item.Kind == "Deployment" {

			var targetState ManagedResourceState
			err := json.Unmarshal([]byte(item.TargetState), &targetState)
			if err != nil {
				log.Println(err.Error())
			}

			var targetImages []string
			for _, container := range targetState.Spec.Template.Spec.Containers {
				targetImages = append(targetImages, container.Image)
			}

			var liveState ManagedResourceState
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

func PrepareEnvironment(item interface{}) codefresh2.Environment {

	converted := item.(*unstructured.Unstructured)

	status := converted.Object["status"].(map[string]interface{})
	spec := converted.Object["spec"].(map[string]interface{})
	source := spec["source"].(map[string]interface{})

	healthStatus := status["health"].(map[string]interface{})
	syncStatusObj := status["sync"].(map[string]interface{})

	syncStatus := syncStatusObj["status"].(string)
	syncRevision := syncStatusObj["revision"].(string)

	metadata := converted.Object["metadata"].(map[string]interface{})

	name := metadata["name"].(string)
	historyList := status["history"].([]interface{})

	historyItem := historyList[len(historyList)-1].(map[string]interface{})

	env := codefresh2.Environment{
		HealthStatus: healthStatus["status"].(string),
		SyncStatus:   syncStatus,
		SyncRevision: syncRevision,
		HistoryId:    historyItem["id"].(int64),
		Name:         name,
		Activities:   prepareEnvironmentActivity(name),
		RepoUrl:      source["repoURL"].(string),
	}

	opStateInterface := status["operationState"]

	if opStateInterface != nil {
		operationState := opStateInterface.(map[string]interface{})
		finishedAtInterface := operationState["finishedAt"]
		if finishedAtInterface != nil {
			env.FinishedAt = finishedAtInterface.(string)
		}

	}

	return env

}

func RetrieveName(item interface{}) string {
	converted := item.(*unstructured.Unstructured)
	metadata := converted.Object["metadata"].(map[string]interface{})
	return metadata["name"].(string)
}

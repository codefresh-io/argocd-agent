package argo

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"log"
)

func initDeploymentsStatuses(applicationName string) map[string]string {
	statuses := make(map[string]string)
	resourceTree, _ := GetResourceTree(applicationName)
	for _, node := range resourceTree.Nodes {
		if node.Kind == "Deployment" {
			if node.Health.status == "" {
				statuses[node.Uid] = "Missing"
			} else {
				statuses[node.Uid] = node.Health.status
			}

		}
	}
	return statuses
}

func prepareEnvironmentActivity(applicationName string) *EnvironmentDeploymentItem {

	resource := GetManagedResources(applicationName)

	statuses := initDeploymentsStatuses(applicationName)

	for _, item := range resource.Items {
		if item.Kind == "Deployment" {

			var targetState ManagedResourceState
			err := json.Unmarshal([]byte(item.TargetState), &targetState)
			if err != nil {
				log.Println(err.Error())
			}

			var liveState ManagedResourceState
			err = json.Unmarshal([]byte(item.LiveState), &liveState)
			if err != nil {
				log.Println(err.Error())
			}

			log.Println(liveState)

			status := statuses[liveState.Metadata.Uid]

			log.Println("Live deployment status " + status)
		}
	}

	return nil
}

func PrepareEnvironment(applicationName string, item interface{}) Environment {
	converted := item.(*unstructured.Unstructured)

	status := converted.Object["status"].(map[string]interface{})

	operationState := status["operationState"].(map[string]interface{})
	finishedAt := operationState["finishedAt"].(string)

	healthStatus := status["health"].(map[string]interface{})

	syncStatusObj := status["sync"].(map[string]interface{})

	syncStatus := syncStatusObj["status"].(string)
	syncRevision := syncStatusObj["revision"].(string)

	metadata := converted.Object["metadata"].(map[string]interface{})
	name := metadata["name"].(string)

	return Environment{
		FinishedAt:   finishedAt,
		HealthStatus: healthStatus["status"].(string),
		SyncStatus:   syncStatus,
		SyncRevision: syncRevision,
		Name:         name,
	}

}

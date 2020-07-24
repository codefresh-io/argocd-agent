package argo

import (
	"encoding/json"
	"github.com/codefresh-io/argocd-listener/src/codefresh"
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

func prepareEnvironmentActivity(applicationName string) []codefresh.EnvironmentActivity {

	resource := GetManagedResources(applicationName)

	statuses := initDeploymentsStatuses(applicationName)

	var activities []codefresh.EnvironmentActivity

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

			log.Println(liveState)

			status := statuses[liveState.Metadata.Uid]

			log.Println("Live deployment status " + status)

			activities = append(activities, codefresh.EnvironmentActivity{
				Name:         item.Name,
				TargetImages: targetImages,
				Status:       status,
				LiveImages:   liveImages,
			})
		}
	}

	return activities
}

func PrepareEnvironment(applicationName string, item interface{}) codefresh.Environment {
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

	return codefresh.Environment{
		FinishedAt:   finishedAt,
		HealthStatus: healthStatus["status"].(string),
		SyncStatus:   syncStatus,
		SyncRevision: syncRevision,
		Name:         name,
		Activities:   prepareEnvironmentActivity(applicationName),
	}

}

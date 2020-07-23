package argo

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog"
)

func process(item interface{}) {

	converted := item.(*unstructured.Unstructured)

	status := converted.Object["status"].(map[string]interface{})

	operationState := status["operationState"].(map[string]interface{})
	finishedAt := operationState["finishedAt"]

	klog.Info(finishedAt)

	healthStatus := status["health"].(map[string]interface{})

	klog.Info(healthStatus["status"])

	syncStatusObj := status["sync"].(map[string]interface{})

	syncStatus := syncStatusObj["status"]
	syncRevision := syncStatusObj["revision"]

	klog.Info(syncStatus)
	klog.Info(syncRevision)

	metadata := converted.Object["metadata"].(map[string]interface{})
	name := metadata["name"]

	klog.Info(name)

	// retrieve additional info
	GetResourceTree("task", "https://34.71.103.174")

}

package argo

import argo "github.com/codefresh-io/argocd-sdk/pkg/api"

type ManagedResourceState struct {
	Spec     ManagedResourceStateSpec
	Metadata ManagedResourceStateMetadata
	Status   ManagedResourceStateStatus
}

type ManagedResourceStateMetadata struct {
	Uid string
}

type ManagedResourceStateSpec struct {
	Replicas int64
	Template ManagedResourceStateTemplate
}

type ManagedResourceStateStatus struct {
	Replicas           int64
	ReadyReplicas      int64
	UpdatedReplicas    int64
	UnavaiableReplicas int64
}

type ManagedResourceStateTemplate struct {
	Spec ManagedResourceTemplateSpec
}

type ManagedResourceTemplateSpec struct {
	Containers []ManagedResourceTemplateContainer
}

type ManagedResourceTemplateContainer struct {
	Image string
}

type ApplicationResource struct {
	Name      string      `json:"name"`
	Kind      string      `json:"kind"`
	Namespace string      `json:"namespace"`
	Status    string      `json:"status"`
	Health    argo.Health `json:"health"`
}

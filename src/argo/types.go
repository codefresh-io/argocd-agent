package argo

type ResourceTree struct {
	Nodes []Node
}

type Node struct {
	Kind   string
	Uid    string
	Health Health
}

type Health struct {
	status string
}

type ManagedResource struct {
	Items []ManagedResourceItem
}

type ManagedResourceItem struct {
	Kind        string
	TargetState string
	LiveState   string
	Name        string
}

type ManagedResourceState struct {
	Spec     ManagedResourceStateSpec
	Metadata ManagedResourceStateMetadata
}

type ManagedResourceStateMetadata struct {
	Uid string
}

type ManagedResourceStateSpec struct {
	Template ManagedResourceStateTemplate
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

type Environment struct {
	FinishedAt   string
	HealthStatus string
	SyncStatus   string
	SyncRevision string
	Name         string
	Activities   []EnvironmentActivity
}

type EnvironmentActivity struct {
	Name         string
	TargetImages []string
	Status       string
	LiveImages   []string
}

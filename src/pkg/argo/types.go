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
	Status string
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

type Project struct {
	Items []ProjectItem
}

type ProjectItem struct {
	Metadata ProjectMetadata
}

type ProjectMetadata struct {
	Name string
}

type Application struct {
	Items []ApplicationItem
}

type ApplicationItem struct {
	Metadata ApplicationMetadata
}

type ApplicationMetadata struct {
	Name string
}

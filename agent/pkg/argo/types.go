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
	Status string `json:"status"`
}

type ManagedResource struct {
	Items []ManagedResourceItem
}

type ServerInfo struct {
	Version string
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

type Project struct {
	Items []ProjectItem
}

type ProjectItem struct {
	Metadata ProjectMetadata `json:"metadata"`
}

type ProjectMetadata struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

type Application struct {
	Items []ApplicationItem
}

type ApplicationItem struct {
	Metadata ApplicationMetadata `json:"metadata"`
	Spec     ApplicationSpec     `json:"spec"`
}

type ApplicationMetadata struct {
	Name        string `json:"name"`
	UID         string `json:"uid"`
	Namespace   string `json:"namespace"`
	ClusterName string `json:"clusterName"`
}

type ApplicationSpecDestination struct {
	Server    string `json:"server"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type ApplicationSpec struct {
	Project     string                     `json:"project"`
	Destination ApplicationSpecDestination `json:"destination"`
}

type ApplicationResource struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Health    Health `json:"health"`
}

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
		Project    string
		SyncPolicy struct {
			Automated interface{}
		}
	}
	Metadata struct {
		Name string
	}
}

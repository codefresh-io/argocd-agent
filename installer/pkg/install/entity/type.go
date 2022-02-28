package entity

type InstallCmdOptions struct {
	Kube Kube
	Argo ArgoOptions
	Env  struct {
		Name string
	}
	NewRelic struct {
		Key string
	}
	Codefresh struct {
		Host                   string
		Token                  string
		Integration            string
		Suffix                 string
		SyncMode               string
		ApplicationsForSync    string
		ApplicationsForSyncArr []string
		Provider               string
	}
	Git struct {
		Integration string
		Password    string
	}
	Host struct {
		HttpProxy  string
		HttpsProxy string
	}
	Agent struct {
		Version string
	}
	Replicas int
}

type Kube struct {
	ManifestPath string
	Namespace    string
	InCluster    bool
	Context      string
	NodeSelector string
	ConfigPath   string

	MasterUrl   string
	BearerToken string
	ClusterName string
}

type ArgoOptions struct {
	Host     string
	Username string
	Password string
	Token    string
	Update   bool
	FailFast bool
}

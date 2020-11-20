package install

type InstallCmdOptions struct {
	Kube      Kube
	Argo      ArgoOptions
	Codefresh struct {
		Host                   string
		Token                  string
		Integration            string
		SyncMode               string
		ApplicationsForSync    string
		ApplicationsForSyncArr []string
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
}

package install

type InstallCmdOptions struct {
	Kube Kube
	Argo struct {
		Host     string
		Username string
		Password string
		Token    string
		Update   bool
	}
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
	Agent struct {
		Version string
	}
}

type Kube struct {
	Namespace    string
	InCluster    bool
	Context      string
	NodeSelector string
	ConfigPath   string

	MasterUrl   string
	BearerToken string
}

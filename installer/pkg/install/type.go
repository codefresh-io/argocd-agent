package install

type InstallCmdOptions struct {
	Kube struct {
		Namespace    string
		InCluster    bool
		Context      string
		NodeSelector string
		ConfigPath   string

		MasterUrl   string
		BearerToken string
	}
	Argo struct {
		Host     string
		Username string
		Password string
		Token    string
	}
	Codefresh struct {
		Host                string
		Token               string
		Integration         string
		SyncMode            string
		ApplicationsForSync string
	}
	Agent struct {
		Version string
	}
}

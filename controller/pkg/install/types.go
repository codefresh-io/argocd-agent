package install

type CmdOptions struct {
	Git struct {
		Auth struct {
			Type string
			Pass string
		}
		Integration string
		RepoUrl     string
	}

	Codefresh struct {
		Host string
		Auth struct {
			Token string
		}
		Clusters []string
	}

	Kube struct {
		ManifestPath string
		Namespace    string
		Context      string
		ConfigPath   string
	}
	Argo struct {
		Password string
	}
}

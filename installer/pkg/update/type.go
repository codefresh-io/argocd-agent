package update

type CmdOptions struct {
	Codefresh struct {
		Suffix string
	}
	Kube struct {
		Namespace  string
		InCluster  bool
		Context    string
		ConfigPath string
	}
}

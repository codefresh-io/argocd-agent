package update

type CmdOptions struct {
	Kube struct {
		Namespace  string
		InCluster  bool
		Context    string
		ConfigPath string
	}
}

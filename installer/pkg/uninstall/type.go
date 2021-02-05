package uninstall

type UninstallCmdOptions struct {
	Kube struct {
		Namespace  string
		InCluster  bool
		Context    string
		ConfigPath string
	}
}

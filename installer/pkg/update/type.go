package update

import "github.com/codefresh-io/argocd-listener/installer/pkg/install"

type CmdOptions struct {
	Agent struct {
		Version string
	}
	Codefresh struct {
		Suffix string
	}
	Kube install.Kube
}

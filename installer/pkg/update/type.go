package update

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/type"
)

type CmdOptions struct {
	Codefresh struct {
		Suffix string
	}
	Kube _type.Kube
}

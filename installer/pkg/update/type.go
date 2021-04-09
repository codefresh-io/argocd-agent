package update

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
)

type CmdOptions struct {
	Codefresh struct {
		Suffix string
	}
	Kube entity.Kube
}

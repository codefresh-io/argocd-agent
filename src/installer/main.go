package installer

import (
	"github.com/codefresh-io/argocd-listener/src/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/src/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/src/installer/cmd"
)

func main() {
	token := argo.GetToken("admin", "newpassword", "https://34.71.103.174")
	store.SetToken(token)

	_ = cmd.Execute()

	//argo.Watch()
}

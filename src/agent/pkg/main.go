package installer

import (
	"github.com/codefresh-io/argocd-listener/src/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/src/agent/pkg/store"
	"os"
)

func main() {
	token := argo.GetToken(os.Getenv("ARGO_HOST"), os.Getenv("ARGO_USERNAME"), os.Getenv("ARGO_PASSWORD"))
	store.SetArgoToken(token)

	store.SetCodefresh(os.Getenv("CODEFRESH_HOST"), os.Getenv("CODEFRESH_TOKEN"))

	argo.Watch()
}

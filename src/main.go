package main

import (
	"github.com/codefresh-io/argocd-listener/src/pkg/argo"
	"github.com/codefresh-io/argocd-listener/src/pkg/store"
)

func main() {
	token := argo.GetToken("admin", "newpassword", "https://34.71.103.174")
	store.SetToken(token)
	argo.Watch()
}

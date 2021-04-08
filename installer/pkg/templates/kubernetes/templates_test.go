package kubernetes

import "testing"

func TestTemplatesMap(t *testing.T) {
	templates := TemplatesMap()
	secrets := templates["4_secret.yaml"]

	originalSecret := `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: cf-argocd-agent{{ .Codefresh.Suffix }}
  namespace: {{ .Namespace }}
data:
  codefresh.token: {{ .Codefresh.Token }}
  {{- if .Argo.Token }}
  argo.token: {{ .Argo.Token }}
  {{- end }}
  {{- if .Kube.BearerToken }}
  kube.bearertoken: {{ .Kube.BearerToken }}
  {{- end }}
  {{- if .Argo.Password }}
  argo.password: {{ .Argo.Password }}
  {{- end }}`

	if secrets != originalSecret {
		t.Error("Original secrets and generated secrets are different")
	}

}

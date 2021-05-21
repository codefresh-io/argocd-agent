package kubernetes

import "testing"

var _ = func() bool {
	testing.Init()
	return true
}()

func TestTemplatesMapSecret(t *testing.T) {
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
  {{- end }}
  { { - if .NewRelic.Key } }
  newrelic.key: { { .NewRelic.Key } }
  { { - end } }`

	if secrets != originalSecret {
		t.Error("Original secrets and generated secrets are different")
	}

}

func TestTemplatesMapSa(t *testing.T) {
	templates := TemplatesMap()
	resource := templates["1_sa.yaml"]

	originalResource := `apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app: cf-argocd-agent{{ .Codefresh.Suffix }}
  name: cf-argocd-agent{{ .Codefresh.Suffix }}
  namespace: {{ .Namespace }}`

	if resource != originalResource {
		t.Error("Original sa and generated sa are different")
	}
}

func TestTemplatesMapClusterRole(t *testing.T) {
	templates := TemplatesMap()
	resource := templates["2_cluster_role.yaml"]

	originalResource := `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    app: cf-argocd-agent{{ .Codefresh.Suffix }}
  name: cf-argocd-agent{{ .Codefresh.Suffix }}
rules:
  - apiGroups:
      - argoproj.io
    resources:
      - applications
      - appprojects
    verbs:
      - get
      - list
      - watch
`

	if resource != originalResource {
		t.Error("Original cluster role and generated cluster role are different")
	}
}

func TestTemplatesMapClusterRoleBinding(t *testing.T) {
	templates := TemplatesMap()
	resource := templates["3_cluster_role_binding.yaml"]

	originalResource := `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    app: cf-argocd-agent{{ .Codefresh.Suffix }}
  name: cf-argocd-agent{{ .Codefresh.Suffix }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cf-argocd-agent{{ .Codefresh.Suffix }}
subjects:
  - kind: ServiceAccount
    name: cf-argocd-agent{{ .Codefresh.Suffix }}
    namespace: {{ .Namespace }}`

	if resource != originalResource {
		t.Error("Original cluster role binding and generated cluster role binding are different")
	}
}

func TestTemplatesMapDeployment(t *testing.T) {
	templates := TemplatesMap()
	resource := templates["5_deployment.yaml"]

	originalResource := `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: cf-argocd-agent{{ .Codefresh.Suffix }}
  name: cf-argocd-agent{{ .Codefresh.Suffix }}
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      app: cf-argocd-agent{{ .Codefresh.Suffix }}
  replicas: 1
  revisionHistoryLimit: 5
  strategy:
    rollingUpdate:
      maxSurge: 50%
      maxUnavailable: 50%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: cf-argocd-agent{{ .Codefresh.Suffix }}
    spec:
      serviceAccountName: cf-argocd-agent{{ .Codefresh.Suffix }}
      containers:
      - env:
        {{- if .Host.HttpProxy }}
        - name: HTTP_PROXY
          value: {{ .Host.HttpProxy }}
        - name: http_proxy
          value: {{ .Host.HttpProxy }}
        {{- end }}
        {{- if .Host.HttpsProxy }}
        - name: HTTPS_PROXY
          value: {{ .Host.HttpsProxy }}
        - name: https_proxy
          value: {{ .Host.HttpsProxy }}
        {{- end }}
        - name: AGENT_VERSION
          value: "{{ .Agent.Version }}"
        - name: ARGO_HOST
          value: {{ .Argo.Host }}
        - name: ARGO_USERNAME
          value: {{ .Argo.Username }}
        {{- if .Argo.Password }}
        - name: ARGO_PASSWORD
          valueFrom:
            secretKeyRef:
              name: cf-argocd-agent{{ .Codefresh.Suffix }}
              key: argo.password
        {{- end }}
        {{- if .Argo.Token }}
        - name: ARGO_TOKEN
          valueFrom:
            secretKeyRef:
              name: cf-argocd-agent{{ .Codefresh.Suffix }}
              key: argo.token
        {{- end }}
        - name: ENV_NAME
          value: {{ .Env.Name }}
        {{ - if .NewRelic.Key }}
        - name: NEWRELIC_LICENSE_KEY
          valueFrom:
            secretKeyRef:
              name: cf-argocd-agent{{ .Codefresh.Suffix }}
              key: newrelic.key
        {{ - end }}
        - name: CODEFRESH_HOST
          value: {{ .Codefresh.Host }}
        - name: CODEFRESH_TOKEN
          valueFrom:
            secretKeyRef:
              name: cf-argocd-agent{{ .Codefresh.Suffix }}
              key: codefresh.token
        - name: IN_CLUSTER
          value: "{{ .Kube.InCluster }}"
        - name: MASTERURL
          value: "{{ .Kube.MasterUrl }}"
        {{- if .Kube.BearerToken }}
        - name: BEARERTOKEN
          valueFrom:
            secretKeyRef:
              name: cf-argocd-agent{{ .Codefresh.Suffix }}
              key: kube.bearertoken
        {{- end }}
        - name: SYNC_MODE
          value: "{{ .Codefresh.SyncMode }}"
        - name: APPLICATIONS_FOR_SYNC
          value: "{{ .Codefresh.ApplicationsForSync }}"
        - name: CODEFRESH_INTEGRATION
          value: {{ .Codefresh.Integration }}
        - name: CODEFRESH_GIT_INTEGRATION
          value: {{ .Git.Integration }}
        image: codefresh/argocd-agent:stable
        imagePullPolicy: Always
        name: cf-argocd-agent{{ .Codefresh.Suffix }}
        resources:
          requests:
            memory: "256Mi"
            cpu: "0.4"
          limits:
            memory: "512Mi"
            cpu: "0.8"
      restartPolicy: Always
      nodeSelector:
        kubernetes.io/arch: amd64
`

	if resource != originalResource {
		t.Error("Original deployment and generated deployment are different")
	}
}

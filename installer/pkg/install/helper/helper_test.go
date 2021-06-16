package helper

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func verifySummaryItem(t *testing.T, expectedTitle string, expectedValue string, item SummaryItem) {
	if item.message != expectedTitle {
		t.Errorf("Wrong message %s, should be %s", item.message, expectedTitle)
	}

	if item.value != expectedValue {
		t.Errorf("Wrong value %s, should be %s", item.value, expectedValue)
	}
}

func TestGetDefaultKubeConfigPath(t *testing.T) {
	expectedValue := "/someDir/.kube/config"
	kubeConfigPath := GetDefaultKubeConfigPath("/someDir")
	if kubeConfigPath != expectedValue {
		t.Errorf("Wrong default kubeconfig path %s, should be %s", kubeConfigPath, expectedValue)
	}
}

func TestBuildSummary(t *testing.T) {

	options := entity.InstallCmdOptions{
		Kube: entity.Kube{
			ClusterName: "ctx",
			Namespace:   "ns",
		},
		Argo: entity.ArgoOptions{
			Host:     "ah",
			Username: "u",
			Password: "p",
		},
		Codefresh: struct {
			Host                   string
			Token                  string
			Integration            string
			Suffix                 string
			SyncMode               string
			ApplicationsForSync    string
			ApplicationsForSyncArr []string
			Provider               string
		}{Host: "localhost"},
		Git: struct {
			Integration string
			Password    string
		}{Integration: "argo"},
		Host: struct {
			HttpProxy  string
			HttpsProxy string
		}{HttpProxy: "hp", HttpsProxy: "hps"},
		Agent: struct {
			Version string
		}{},
	}

	summaryItems := buildSummary(&options)

	if len(summaryItems) != 10 {
		t.Error("Wrong amount of summary items")
	}

	verifySummaryItem(t, "Kubernetes Context", "ctx", summaryItems[0])
	verifySummaryItem(t, "Kubernetes Namespace", "ns", summaryItems[1])
	verifySummaryItem(t, "Git Integration", "argo", summaryItems[2])
	verifySummaryItem(t, "Codefresh Host", "localhost", summaryItems[3])
	verifySummaryItem(t, "ArgoCD Host", "ah", summaryItems[4])
	verifySummaryItem(t, "ArgoCD Username", "u", summaryItems[5])
	verifySummaryItem(t, "ArgoCD Password", "******", summaryItems[6])
	verifySummaryItem(t, "Enable auto-sync of applications", "No", summaryItems[7])
	verifySummaryItem(t, "HTTP proxy", "hp", summaryItems[8])
	verifySummaryItem(t, "HTTPS proxy", "hps", summaryItems[9])

}

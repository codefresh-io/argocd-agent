package helper

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
)

func ShowSummary(installOptions *install.InstallCmdOptions) {
	logger.Success("\nInstallation options summary:")
	var items []SummaryItem
	var syncModeStr string
	if installOptions.Codefresh.SyncMode == codefresh.ContinueSync {
		syncModeStr = "Yes"
	} else {
		syncModeStr = "No"
	}

	items = append(items, SummaryItem{
		message: "Kubernetes Context", value: installOptions.Kube.ClusterName,
	})
	items = append(items, SummaryItem{
		message: "Kubernetes Namespace",
		value:   installOptions.Kube.Namespace,
	})
	items = append(items, SummaryItem{
		message: "Git Integration",
		value:   installOptions.Git.Integration,
	})
	items = append(items, SummaryItem{
		message: "Codefresh Host",
		value:   installOptions.Codefresh.Host,
	})
	items = append(items, SummaryItem{
		message: "ArgoCD Host",
		value:   installOptions.Argo.Host,
	})
	if installOptions.Argo.Password != "" {
		items = append(items, SummaryItem{
			message: "ArgoCD Username",
			value:   installOptions.Argo.Username,
		})
		items = append(items, SummaryItem{
			message: "ArgoCD Password",
			value:   "******",
		})
	} else if installOptions.Argo.Token != "" {
		items = append(items, SummaryItem{
			message: "ArgoCD Token",
			value:   installOptions.Argo.Token,
		})
	}

	items = append(items, SummaryItem{
		message: "Enable auto-sync of applications",
		value:   syncModeStr,
	})
	items = append(items, SummaryItem{
		message: "HTTP proxy",
		value:   getProxyString(installOptions.Host.HttpProxy),
	})
	items = append(items, SummaryItem{
		message: "HTTPS proxy",
		value:   getProxyString(installOptions.Host.HttpsProxy),
	})

	for i, item := range items {
		logger.Summary(i+1, item.message, item.value)
	}
	logger.Info("")
}

func getProxyString(proxyValue string) string {
	if proxyValue != "" {
		return proxyValue
	}
	return "none"

}

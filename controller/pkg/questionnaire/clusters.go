package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/controller/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
)

func AskAboutClusters(installOptions *install.CmdOptions, clusters []*codefresh.ClusterMinified) error {
	if len(clusters) < 1 {
		return nil
	}

	clustersSelectors := make([]string, 0)
	for _, cluster := range clusters {
		clustersSelectors = append(clustersSelectors, cluster.Selector)
	}

	_, clustersForSync := prompt.Multiselect(clustersSelectors, "Select clusters your would be like to register")
	installOptions.Codefresh.Clusters = clustersForSync
	return nil
}

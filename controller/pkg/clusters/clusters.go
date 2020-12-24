package clusters

import (
	"encoding/base64"
	"fmt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	argo "github.com/codefresh-io/argocd-sdk/pkg/api"
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
)

const CODEFRESH_CLUSTER_PREFIX = "cf-"

func filterClusters(clusters []*codefresh.ClusterMinified) []*codefresh.ClusterMinified {
	filteredClusters := []*codefresh.ClusterMinified{}
	for _, cluster := range clusters {
		if cluster.Provider == "local" && cluster.BehindFirewall == false {
			filteredClusters = append(filteredClusters, cluster)
		}
	}
	return filteredClusters
}

func GetAvailableClusters(cfClustersApi codefresh.IClusterAPI) ([]*codefresh.ClusterMinified, error) {
	clustersList, err := cfClustersApi.GetAccountClusters()
	if err != nil {
		return []*codefresh.ClusterMinified{}, err
	}
	return filterClusters(clustersList), nil
}

func ImportFromCodefresh(clusters []string, cfClustersApi codefresh.IClusterAPI, argoClustersApi argo.ClusterApi) error {
	if len(clusters) < 1 {
		logger.Warning(fmt.Sprint("Import clusters skipped because nothing was selected..."))
		return nil
	}
	logger.Info(fmt.Sprint("Import clusters..."))
	for _, clusterSelector := range clusters {
		cluster, err := cfClustersApi.GetClusterCredentialsByAccountId(clusterSelector)
		if err != nil {
			return err
		}

		bearer, err := base64.StdEncoding.DecodeString(cluster.Auth.Bearer)
		if err != nil {
			return err
		}

		_, err = argoClustersApi.CreateCluster(argo.ClusterOpt{
			Name:   CODEFRESH_CLUSTER_PREFIX + clusterSelector,
			Server: cluster.Url,
			Config: argo.ClusterConfig{
				BearerToken: string(bearer),
				TlsClientConfig: argo.TlsClientConfig{
					CaData:   cluster.Ca,
					Insecure: false,
				},
			},
		})
		if err != nil {
			return err
		}
		logger.Success(fmt.Sprintf("Successfull created cluster \"%s\"", clusterSelector))
	}

	return nil
}

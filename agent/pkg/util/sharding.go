package util

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"github.com/mitchellh/mapstructure"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Sharding struct {
	numberOfShard int
	replicas      int
	applications  []*argoSdk.ApplicationItem
}

func NewSharding(numberOfShard, replicas int) *Sharding {
	return &Sharding{
		numberOfShard: numberOfShard,
		replicas:      replicas,
		applications:  nil,
	}
}

func (sh *Sharding) applicationsRange(amountOfApps int) (int, int) {
	if amountOfApps == 0 {
		return 0, 0
	}
	// primary := sh.replicas % amountOfApps
	appsPerShard := amountOfApps / sh.replicas
	from := appsPerShard * sh.numberOfShard
	return from, from + appsPerShard
}

func (sh *Sharding) InitApplications(applications []unstructured.Unstructured) {
	from, to := sh.applicationsRange(len(applications))
	pickedApps := applications[from:to]

	var appsForCurrentShard []*argoSdk.ApplicationItem

	for _, app := range pickedApps {
		var appForCurrentShard argoSdk.ApplicationItem
		Convert(app.Object, &appForCurrentShard)
		appsForCurrentShard = append(appsForCurrentShard, &appForCurrentShard)
	}

	if appsForCurrentShard != nil && len(appsForCurrentShard) > 0 {
		for i := 0; i < len(appsForCurrentShard); i++ {
			logger.GetLogger().Infof("[Sharding] Choose \"%s\" for processing", appsForCurrentShard[i].Metadata.Name)
		}
	}

	sh.applications = appsForCurrentShard
}

func (sh *Sharding) ShouldBeProcessed(obj interface{}) bool {
	if sh.applications != nil && len(sh.applications) > 0 {
		var app argoSdk.ApplicationItem
		err := mapstructure.Decode(obj.(*unstructured.Unstructured).Object, &app)
		if err != nil {
			logger.GetLogger().Infof("Failed to parse app , reason %s", err.Error())
		}
		found := false
		for i := 0; i < len(sh.applications); i++ {
			appFromShard := sh.applications[i]
			if appFromShard.Metadata.Name == app.Metadata.Name && appFromShard.Metadata.Namespace == app.Metadata.Namespace {
				found = true
				break
			}
		}

		if !found {
			logger.GetLogger().Infof("Skip handling of \"%s\" app because it is not in this shard", app.Metadata.Name)
		}

		return found
	}

	return true
}

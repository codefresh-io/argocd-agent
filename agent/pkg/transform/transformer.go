package transform

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/git/provider"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"github.com/mitchellh/mapstructure"
	"sort"
)

type EnvTransformer struct {
	argoApi argo.ArgoApi
}

var envTransformer *EnvTransformer

func GetEnvTransformerInstance(argoApi argo.ArgoApi) *EnvTransformer {
	if envTransformer != nil {
		return envTransformer
	}
	envTransformer = &EnvTransformer{
		argoApi,
	}
	return envTransformer
}

func (envTransformer *EnvTransformer) initDeploymentsStatuses(applicationName string) (map[string]string, error) {
	statuses := make(map[string]string)
	resourceTree, err := envTransformer.argoApi.GetResourceTree(applicationName)
	if err != nil {
		return nil, err
	}
	for _, node := range resourceTree.Nodes {
		if node.Health.Status == "" {
			statuses[node.Uid] = "Missing"
		} else {
			statuses[node.Uid] = node.Health.Status
		}
	}
	return statuses, nil
}

func (envTransformer *EnvTransformer) prepareEnvironmentActivity(applicationName string) ([]codefresh2.EnvironmentActivity, error) {

	resource, err := envTransformer.argoApi.GetManagedResources(applicationName)
	if err != nil {
		return nil, err
	}

	statuses, err := envTransformer.initDeploymentsStatuses(applicationName)

	if err != nil {
		return nil, err
	}

	var services = make(map[string]codefresh2.EnvironmentActivity)

	for _, item := range resource.Items {
		var liveState argo.ManagedResourceState
		err = json.Unmarshal([]byte(item.LiveState), &liveState)
		if err != nil {
			logger.GetLogger().Errorf("Failed to unmarshal \"LiveState\" to ManagedResourceState, reason %v", err)
			continue
		}

		var liveImages []string
		for _, container := range liveState.Spec.Template.Spec.Containers {
			if container.Image != "" {
				liveImages = append(liveImages, container.Image)
			}

		}
		if len(liveImages) != 0 {
			status := statuses[liveState.Metadata.Uid]

			replicasStatus := liveState.Status

			fromReplicaState := codefresh2.ReplicaState{
				Current: replicasStatus.ReadyReplicas - (replicasStatus.UpdatedReplicas - replicasStatus.UnavaiableReplicas),
			}

			toReplicasState := codefresh2.ReplicaState{
				Current: replicasStatus.UpdatedReplicas,
				Desired: liveState.Spec.Replicas,
			}

			services[item.Name] = codefresh2.EnvironmentActivity{
				Name:       item.Name,
				Status:     status,
				LiveImages: liveImages,
				ReplicaSet: codefresh2.EnvironmentActivityRS{
					From: fromReplicaState,
					To:   toReplicasState,
				},
			}
		}

	}

	var result = make([]codefresh2.EnvironmentActivity, 0, len(services))

	for _, svc := range services {
		result = append(result, svc)
	}

	return result, nil
}

func filterResources(resources interface{}) []interface{} {
	result := make([]interface{}, 0)
	if resources == nil {
		return result
	}
	for _, resource := range resources.([]interface{}) {
		resourceItem := resource.(map[string]interface{})
		resourceKind := resourceItem["kind"]
		if resourceKind == "Service" || resourceKind == "Pod" {
			result = append(result, resourceItem)
		}
	}
	return result
}

func (envTransformer *EnvTransformer) PrepareEnvironment(envItem map[string]interface{}) (error, *codefresh2.Environment) {

	var app argoSdk.ArgoApplication
	err := mapstructure.Decode(envItem, &app)
	if err != nil {
		return err, nil
	}

	github := provider.NewGithubProvider()

	name := app.Metadata.Name
	historyList := app.Status.History
	revision := app.Status.OperationState.SyncResult.Revision
	repoUrl := app.Spec.Source.RepoURL

	if revision == "" {
		return errors.New("revision is empty"), nil
	}

	resources, err := envTransformer.argoApi.GetResourceTreeAll(name)
	if err != nil {
		return err, nil
	}

	// we still need send env , even if we have problem with retrieve gitops info
	err, gitops := github.GetManifestRepoInfo(repoUrl, revision)

	if err != nil {
		logger.GetLogger().Errorf("Failed to retrieve manifest repo git information , reason: %v", err)
	}

	err, historyId := resolveHistoryId(historyList, app.Status.OperationState.SyncResult.Revision, name)

	if err != nil {
		return err, nil
	}

	activities, err := envTransformer.prepareEnvironmentActivity(name)
	if err != nil {
		return err, nil
	}

	syncPolicy := codefresh2.SyncPolicy{AutoSync: &app.Spec.SyncPolicy != nil && app.Spec.SyncPolicy.Automated != nil}

	env := codefresh2.Environment{
		HealthStatus: app.Status.Health.Status,
		SyncStatus:   app.Status.Sync.Status,
		SyncRevision: revision,
		Gitops:       *gitops,
		HistoryId:    historyId,
		Name:         name,
		Activities:   activities,
		Resources:    filterResources(resources),
		RepoUrl:      repoUrl,
		FinishedAt:   app.Status.OperationState.FinishedAt,
		SyncPolicy:   syncPolicy,
		Date:         app.Status.OperationState.FinishedAt,
	}

	err, commit := github.GetCommitByRevision(repoUrl, revision)

	if commit != nil {
		logger.GetLogger().Infof("Retrieve commit message \"%s\" for repo \"%s\" ", *commit.Message, repoUrl)
		env.Commit = *commit
	}

	return nil, &env

}

func resolveHistoryId(historyList []argoSdk.ApplicationHistoryItem, revision string, name string) (error, int64) {
	if historyList == nil {
		logger.GetLogger().Errorf("can`t find history id for application %s, because history list is empty", name)
		return nil, -1
	}

	sort.Slice(historyList, func(i, j int) bool {
		return historyList[i].Id > historyList[j].Id
	})

	for _, item := range historyList {
		if item.Revision == revision {
			return nil, item.Id
		}
	}
	return fmt.Errorf("can`t find history id for application %s", name), 0
}

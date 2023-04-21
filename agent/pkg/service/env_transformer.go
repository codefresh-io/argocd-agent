package service

import (
	"encoding/json"

	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/git/provider"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type EnvTransformer struct {
	argoApi argo.ArgoAPI
}

// This one should not exist and will be refactored soon, Env transformer should be migrated to service layer
type ETransformer interface {
	PrepareEnvironment(app argoSdk.ArgoApplication, historyId int64) (error, *EnvironmentWrapper)
}

var envTransformer *EnvTransformer

func GetEnvTransformerInstance(argoApi argo.ArgoAPI) ETransformer {
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

func (envTransformer *EnvTransformer) prepareEnvironmentActivity(applicationName string) ([]codefreshSdk.EnvironmentActivity, error) {

	resource, err := envTransformer.argoApi.GetManagedResources(applicationName)
	if err != nil {
		return nil, err
	}

	statuses, err := envTransformer.initDeploymentsStatuses(applicationName)

	if err != nil {
		return nil, err
	}

	var services = make(map[string]codefreshSdk.EnvironmentActivity)

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

			fromReplicaState := codefreshSdk.ReplicaState{
				Current: replicasStatus.ReadyReplicas - (replicasStatus.UpdatedReplicas - replicasStatus.UnavaiableReplicas),
			}

			toReplicasState := codefreshSdk.ReplicaState{
				Current: replicasStatus.UpdatedReplicas,
				Desired: liveState.Spec.Replicas,
			}

			services[item.Name] = codefreshSdk.EnvironmentActivity{
				Name:       item.Name,
				Status:     status,
				LiveImages: liveImages,
				ReplicaSet: codefreshSdk.EnvironmentActivityRS{
					From: fromReplicaState,
					To:   toReplicasState,
				},
			}
		}

	}

	var result = make([]codefreshSdk.EnvironmentActivity, 0, len(services))

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
		if resourceKind == "Service" || resourceKind == "Pod" || resourceKind == "Application" {
			result = append(result, resourceItem)
		}
	}
	return result
}

func (envTransformer *EnvTransformer) PrepareEnvironment(app argoSdk.ArgoApplication, historyId int64) (error, *EnvironmentWrapper) {

	git := provider.GetGitProvider()

	name := app.Metadata.Name
	revision := app.Status.OperationState.SyncResult.Revision
	repoUrl := app.Spec.Source.RepoURL
	parentApp, _ := app.Metadata.Labels["app.kubernetes.io/instance"]

	resources, err := envTransformer.argoApi.GetResourceTreeAll(name)

	logPayloadBefore, _ := json.Marshal(resources)
	logger.GetLogger().Debugf("==> Resources BEFORE filter \"%s\", payload: \n%s", name, logPayloadBefore)

	if err != nil {
		return err, nil
	}
	filteredResources := filterResources(resources)

	logPayloadAfter, _ := json.Marshal(filteredResources)
	logger.GetLogger().Debugf("==> Resources AFTER filter \"%s\", payload: \n%s", name, logPayloadAfter)

	// we still need send env , even if we have problem with retrieve gitops info
	err, gitops := git.GetManifestRepoInfo(repoUrl, revision)

	if err != nil {
		logger.GetLogger().Errorf("Failed to retrieve manifest repo git information , reason: %v", err)
	}

	activities, err := envTransformer.prepareEnvironmentActivity(name)
	if err != nil {
		return err, nil
	}

	syncPolicy := codefreshSdk.SyncPolicy{AutoSync: &app.Spec.SyncPolicy != nil && app.Spec.SyncPolicy.Automated != nil}

	env := codefreshSdk.Environment{
		HealthStatus: app.Status.Health.Status,
		SyncStatus:   app.Status.Sync.Status,
		ParentApp:    parentApp,
		SyncRevision: revision,
		Gitops:       *gitops,
		HistoryId:    historyId,
		Name:         name,
		Activities:   activities,
		Resources:    filteredResources,
		RepoUrl:      repoUrl,
		Namespace:    app.Spec.Destination.Namespace,
		FinishedAt:   app.Status.OperationState.FinishedAt,
		SyncPolicy:   syncPolicy,
		Date:         app.Status.OperationState.FinishedAt,
		Context:      &store.GetStore().Codefresh.Integration,
	}

	err, commit := git.GetCommitByRevision(repoUrl, revision)

	if err != nil {
		logger.GetLogger().Errorf("Failed to retrieve commit message, reason %v", err)
	}

	if commit != nil {
		logger.GetLogger().Infof("Retrieve commit message \"%s\" for repo \"%s\" ", *commit.Message, repoUrl)
		env.Commit = codefreshSdk.Commit{
			Time:    commit.Time,
			Message: commit.Message,
			Avatar:  commit.Avatar,
		}
	}

	envWrapper := &EnvironmentWrapper{
		Environment: env,
	}

	if commit != nil {
		envWrapper.Commit = *commit
	}

	return nil, envWrapper
}

package comparator

import (
	"fmt"
	"reflect"

	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/ulule/deepcopier"
)

type Comparator interface {
	Compare(obj1 interface{}, obj2 interface{}) bool
}

type EnvComparator struct {
}

func compareServices(services1 []codefreshSdk.EnvironmentActivity, services2 []codefreshSdk.EnvironmentActivity) bool {

	for _, svc := range services1 {

		foundService := false

		for _, svc2 := range services2 {
			if svc.Name == svc2.Name {
				foundService = reflect.DeepEqual(svc, svc2)
			}
		}

		if !foundService {
			return false
		}
	}

	return true
}

func (comparator EnvComparator) Compare(obj1 interface{}, obj2 interface{}) bool {

	if obj1 == nil || obj2 == nil {
		return false
	}

	env1 := obj1.(*codefreshSdk.Environment)
	env2 := obj2.(*codefreshSdk.Environment)

	newEnv1 := &codefreshSdk.Environment{}
	newEnv2 := &codefreshSdk.Environment{}

	_ = deepcopier.Copy(env1).To(newEnv1)
	_ = deepcopier.Copy(env2).To(newEnv2)

	newEnv1.Resources = nil
	newEnv2.Resources = nil

	// ignore git metadata itself, we need compare only revisions
	newEnv1.Gitops = codefreshSdk.Gitops{}
	newEnv2.Gitops = codefreshSdk.Gitops{}

	// ignore date because we can rely on status
	newEnv1.Date = ""
	newEnv2.Date = ""

	newEnv1.FinishedAt = ""
	newEnv2.FinishedAt = ""

	sameServices := compareServices(newEnv1.Activities, newEnv2.Activities)

	newEnv1.Activities = nil
	newEnv2.Activities = nil

	// add printdiff

	return reflect.DeepEqual(newEnv1, newEnv2) && sameServices
}

type ArgoAppComparator struct {
}

func (comparator ArgoAppComparator) Compare(obj1 interface{}, obj2 interface{}) bool {

	if obj1 == nil || obj2 == nil {
		return false
	}

	app1Str, err := comparator.buildStateStrFromArgoApp(obj1)
	if err != nil {
		logger.GetLogger().Errorf(err.Error())
		return false
	}
	app2Str, err := comparator.buildStateStrFromArgoApp(obj2)
	if err != nil {
		logger.GetLogger().Errorf(err.Error())
		return false
	}

	equal := app1Str == app2Str
	if !equal {
		logger.GetLogger().Debugf("argo app change detected")
		logger.GetLogger().Diff(app1Str, app2Str)
	}
	return equal
}

func (comparator ArgoAppComparator) buildStateStrFromArgoApp(obj interface{}) (string, error) {
	argoApp, ok := obj.(*argoSdk.ArgoApplication)
	if !ok {
		return "", fmt.Errorf("failed to cast compared object to %T", argoSdk.ArgoApplication{})
	}

	err, historyId := service.NewArgoResourceService().ResolveHistoryId(argoApp.Status.History, argoApp.Status.OperationState.SyncResult.Revision, argoApp.Metadata.Name)
	if err != nil {
		return "", err
	}

	stateId := fmt.Sprintf("%v-%v-%v-%v-%v",
		argoApp.Metadata.Name, historyId, argoApp.Status.Health.Status, argoApp.Status.Sync.Status, argoApp.Status.Sync.Revision)

	logger.GetLogger().Debugf("app: %v, state_id: %v", argoApp.Metadata.Name, stateId)
	return stateId, nil
}

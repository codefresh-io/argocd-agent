package comparator

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/ulule/deepcopier"
	"reflect"
)

type Comparator interface {
	Compare(obj1 interface{}, obj2 interface{}) bool
}

type EnvComparator struct {
}

func compareServices(services1 []codefresh.EnvironmentActivity, services2 []codefresh.EnvironmentActivity) bool {

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

	env1 := obj1.(*codefresh.Environment)
	env2 := obj2.(*codefresh.Environment)

	newEnv1 := &codefresh.Environment{}
	newEnv2 := &codefresh.Environment{}

	_ = deepcopier.Copy(env1).To(newEnv1)
	_ = deepcopier.Copy(env2).To(newEnv2)

	newEnv1.Resources = nil
	newEnv2.Resources = nil

	sameServices := compareServices(newEnv1.Activities, newEnv2.Activities)

	newEnv1.Activities = nil
	newEnv2.Activities = nil

	return reflect.DeepEqual(newEnv1, newEnv2) && sameServices
}

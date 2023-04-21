package util

import (
	"reflect"

	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
)

var previousState = make(map[string]interface{})

func ProcessDataWithFilter(itemType string, key *string, data interface{}, comparator func(oldItem interface{}, newItem interface{}) bool, callback func() error) error {
	stateKey := itemType

	if key != nil {
		stateKey += "." + *key
	}

	oldItem := previousState[stateKey]

	if comparator == nil {
		// default comparator
		comparator = reflect.DeepEqual
	}

	if comparator(oldItem, data) {
		logger.GetLogger().Infof("Item with key \"%s\" didn't change, ignoring the callback", stateKey)
		return nil
	}
	logger.GetLogger().Infof("Item with key \"%s\" was changed, executing the callback", stateKey)

	err := callback()

	if err != nil {
		return err
	}
	previousState[stateKey] = data

	return nil
}

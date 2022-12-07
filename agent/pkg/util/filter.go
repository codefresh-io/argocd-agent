package util

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"

	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var previousState = make(map[string]interface{})

func printDiff(stateKey string, oldItem interface{}, newItem interface{}) error {
	printResourceDiff, _ := os.LookupEnv("PRINT_RESOURCE_DIFF")

	printResourceDiffBool, err := strconv.ParseBool(printResourceDiff)
	if err != nil {
		printResourceDiffBool = false
	}

	if !printResourceDiffBool {
		return nil
	}

	if oldItem == nil || newItem == nil {
		logger.GetLogger().Infof("Ignore diff view because one of the entities is nil")
		return nil
	}

	prevState, err := json.Marshal(oldItem)
	if err != nil {
		return err
	}
	newState, err := json.Marshal(newItem)
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(prevState), string(newState), false)
	logger.GetLogger().Infof(dmp.DiffPrettyText(diffs))

	return nil
}

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
	printDiff(stateKey, oldItem, data)

	err := callback()

	if err != nil {
		return err
	}
	previousState[stateKey] = data

	return nil
}

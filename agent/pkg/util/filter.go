package util

import "reflect"

var previousState = make(map[string]interface{})

func ProcessDataWithFilter(itemType string, data interface{}, callback func() error) error {
	oldItem := previousState[itemType]

	if reflect.DeepEqual(oldItem, data) {
		return nil
	}

	err := callback()

	if err != nil {
		return err
	}
	previousState[itemType] = data

	return nil
}

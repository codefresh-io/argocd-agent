package util

import "reflect"

var previousState = make(map[string]interface{})

func ProcessDataWithFilter(itemType string, data interface{}, callback func() error) (error, interface{}) {
	oldItem := previousState[itemType]

	if reflect.DeepEqual(oldItem, data) {
		return nil, nil
	}

	err := callback()

	if err != nil {
		return err, nil
	}
	previousState[itemType] = data

	return nil, data
}

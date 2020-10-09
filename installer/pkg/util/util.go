package util

func GetMapKeys(obj map[string]string) []string {
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	return keys
}

func ConvertIntToStringArray(entities []interface{}) []string {
	res := make([]string, len(entities))
	for i, v := range entities {
		res[i] = v.(string)
	}
	return res
}

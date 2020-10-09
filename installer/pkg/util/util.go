package util

import "github.com/codefresh-io/argocd-listener/installer/pkg/fs"

func ResolvePackageVersion(versionFromLdFlag string) string {
	// Getting version from ldflag
	if versionFromLdFlag != "" {
		return versionFromLdFlag
	}
	// Getting version from file (for local development)
	return fs.GetPackageVersionFromFile("./VERSION")
}

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

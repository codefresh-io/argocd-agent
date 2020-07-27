package argo

func Allow(name string) bool {
	allowedNames := map[string]bool{
		"task": true,
	}
	return allowedNames[name]
}

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

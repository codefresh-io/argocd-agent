package fs

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"io/ioutil"
	"strings"
)

func ReadFile(pathToFile string) (error, string) {
	content, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return err, ""
	}
	return nil, string(content)
}

func GetPackageVersionFromFile(pathToFile string) string {
	content, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		logger.GetLogger().Errorf("Failed to retrieve package version, reason %v", err)
		return ""
	} else {
		version := getVersionFromContentString(string(content))
		return version
	}
}

func getVersionFromContentString(content string) string {
	return strings.TrimFunc(content, func(r rune) bool {
		_r := string(r)
		return _r == " " || _r == "\n" || _r == "\t" || _r == "\r"
	})
}

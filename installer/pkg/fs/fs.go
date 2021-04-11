package fs

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"io/ioutil"
	"os"
	"strings"
)

func WriteFile(pathToFile string, content string) error {
	f, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
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

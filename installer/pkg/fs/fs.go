package fs

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func GetAgentVersion(pathToFIle string) string {
	content, err := ioutil.ReadFile(pathToFIle)
	if err != nil {
		fmt.Errorf(err.Error())
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

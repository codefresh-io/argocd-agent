package fs

import (
	"io/ioutil"
	"regexp"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestGetPackageVersionFromFile(t *testing.T) {
	content, _ := ioutil.ReadFile("../../VERSION")
	versionFromFile := string(content)
	trimVersion := GetPackageVersionFromFile("../../VERSION")
	re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	if versionFromFile != trimVersion || versionFromFile == "" {
		t.Errorf("Agent version is not valid! expected: >>%v<<, got: >>%v<<", trimVersion, versionFromFile)
	} else if !re.Match([]byte(versionFromFile)) {
		t.Errorf("Agent version is not in valid format! got: >>%v<<", versionFromFile)
	}
}

func TestGetVersionFromContentString(t *testing.T) {
	expectedVersion := "0.1.2"
	testPayload := []string{
		"0.1.2",
		" 0.1.2  ",
		" 0.1.2",
		"0.1.2   ",
		"\t0.1.2   \n",
		"\t0.1.2\n",
		"\t0.1.2",
	}

	for _, payload := range testPayload {
		version := getVersionFromContentString(payload)
		if version != expectedVersion {
			t.Errorf("'GetVersionFromContentString' check version failed, expected: %v, got: %v", expectedVersion, version)
		}
	}
}

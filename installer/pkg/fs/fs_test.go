package fs

import (
	"io/ioutil"
	"testing"
)

func TestGetAgentVersion(t *testing.T) {
	content, _ := ioutil.ReadFile("../../../agent/VERSION")
	versionFromFile := string(content)
	trimVersion := GetAgentVersion("../../../agent/VERSION")
	if versionFromFile != trimVersion || versionFromFile == "" {
		t.Errorf("Agent version is not valid! expected: >>%v<<, got: >>%v<<", trimVersion, versionFromFile)
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

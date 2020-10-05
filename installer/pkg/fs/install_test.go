package fs

import (
	"testing"
)

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
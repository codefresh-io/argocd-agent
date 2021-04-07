package startup

import (
	"testing"
)

func TestValidInput(t *testing.T) {
	input := &Input{
		argoHost:                 "http://argo-host",
		argoToken:                "1234",
		argoUsername:             "",
		argoPassword:             "",
		codefreshToken:           "token",
		codefreshHost:            "http://cf-host",
		codefreshIntegrationName: "test",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}

	err := NewInputValidator(input).Validate()

	if err != nil {
		t.Error("Validation should be passed, input is correct")
	}
}

func testFailedValidation(t *testing.T, input *Input, message string) {
	err := NewInputValidator(input).Validate()

	if err == nil || err.Error() != message {
		t.Errorf("Validation should be failed with message \"%s\"", message)
	}
}

func TestInvalidArgoHost(t *testing.T) {
	input := &Input{
		argoHost:                 "",
		argoToken:                "1234",
		argoUsername:             "",
		argoPassword:             "",
		codefreshToken:           "token",
		codefreshHost:            "http://cf-host",
		codefreshIntegrationName: "test",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}
	testFailedValidation(t, input, "ARGO_HOST variable doesnt exist")
}

func TestInvalidArgoUsername(t *testing.T) {
	input := &Input{
		argoHost:                 "http://host",
		argoToken:                "",
		argoUsername:             "",
		argoPassword:             "",
		codefreshToken:           "token",
		codefreshHost:            "http://cf-host",
		codefreshIntegrationName: "test",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}
	testFailedValidation(t, input, "ARGO_USERNAME variable doesnt exist")
}

func TestInvalidArgoPassword(t *testing.T) {
	input := &Input{
		argoHost:                 "http://host",
		argoToken:                "",
		argoUsername:             "1234",
		argoPassword:             "",
		codefreshToken:           "token",
		codefreshHost:            "http://cf-host",
		codefreshIntegrationName: "test",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}
	testFailedValidation(t, input, "ARGO_PASSWORD variable doesnt exist")
}

func TestInvalidCodefreshToken(t *testing.T) {
	input := &Input{
		argoHost:                 "http://host",
		argoToken:                "",
		argoUsername:             "1234",
		argoPassword:             "1234",
		codefreshToken:           "",
		codefreshHost:            "http://cf-host",
		codefreshIntegrationName: "test",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}
	testFailedValidation(t, input, "CODEFRESH_TOKEN variable doesnt exist")
}

func TestInvalidCodefreshIntegration(t *testing.T) {
	input := &Input{
		argoHost:                 "http://host",
		argoToken:                "",
		argoUsername:             "1234",
		argoPassword:             "1234",
		codefreshToken:           "1234",
		codefreshHost:            "http://cf-host",
		codefreshIntegrationName: "",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}
	testFailedValidation(t, input, "CODEFRESH_INTEGRATION variable doesnt exist")
}

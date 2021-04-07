package startup

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"os"
)

type InputValidator struct {
	input *Input
}

func NewInputValidator(input *Input) *InputValidator {
	return &InputValidator{input}
}

func (validator *InputValidator) Validate() error {
	input := validator.input

	if input.argoHost == "" {
		return errors.New("ARGO_HOST variable doesnt exist")
	}

	if input.argoToken == "" {

		if input.argoUsername == "" {
			return errors.New("ARGO_USERNAME variable doesnt exist")
		}

		if input.argoPassword == "" {
			return errors.New("ARGO_PASSWORD variable doesnt exist")
		}

	}

	if input.codefreshToken == "" {
		return errors.New("CODEFRESH_TOKEN variable doesnt exist")
	}

	if input.codefreshIntegrationName == "" {
		return errors.New("CODEFRESH_INTEGRATION variable doesnt exist")
	}

	_, gitIntegrationExistence := os.LookupEnv("CODEFRESH_GIT_INTEGRATION")

	if !gitIntegrationExistence {
		logger.GetLogger().Errorf("No git context")
	}

	return nil
}

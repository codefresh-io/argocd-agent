package acceptance

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
)

type (
	acceptanceTest interface {
		check(argoOptions *entity.ArgoOptions) error
		getMessage() string
	}

	IAcceptanceTestRunner interface {
		Verify(argoOptions *entity.ArgoOptions) error
	}

	AcceptanceTestRunner struct {
	}
)

var tests []acceptanceTest
var runner IAcceptanceTestRunner

// New create runner and init tests suite
func New() IAcceptanceTestRunner {
	if runner == nil {
		// should be first in tests array because we setup token to storage , it is super not good and should be rewritten
		tests = append(tests, &ArgoCredentialsAcceptanceTest{})

		tests = append(tests, &ProjectAcceptanceTest{})
		tests = append(tests, &ApplicationAcceptanceTest{})

		runner = AcceptanceTestRunner{}
	}
	return runner
}

// Verify execute test suites with specific options
func (runner AcceptanceTestRunner) Verify(argoOptions *entity.ArgoOptions) error {
	logger.Info("\nTesting requirements")
	logger.Info("--------------------")
	defer logger.Info("--------------------\n")

	var err error

	for _, test := range tests {
		err = test.check(argoOptions)
		if err != nil {
			logger.FailureTest(test.getMessage())
			return err
		}
		logger.SuccessTest(test.getMessage())
	}

	return nil
}

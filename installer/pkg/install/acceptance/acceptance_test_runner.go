package acceptance

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

type (
	acceptanceTest interface {
		check(argoOptions *entity.ArgoOptions) (error, bool)
		getMessage() string
		failure(argoOptions *entity.ArgoOptions) bool
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
		tests = append(tests, &ArgoAccessibilityAcceptanceTest{prompt: prompt.NewPrompt()})

		// should be before other tests array because we setup token to storage , it is super not good and should be rewritten
		tests = append(tests, &ArgoCredentialsAcceptanceTest{})

		tests = append(tests, &ProjectAcceptanceTest{})
		tests = append(tests, &ApplicationAcceptanceTest{prompt: prompt.NewPrompt()})

		runner = AcceptanceTestRunner{}
	}
	return runner
}

// Verify execute test suites with specific options
func (runner AcceptanceTestRunner) Verify(argoOptions *entity.ArgoOptions) error {
	logger.Info("\nTesting requirements")
	logger.Info("--------------------")
	defer logger.Info("--------------------\n")

	for _, test := range tests {
		err, failfast := test.check(argoOptions)
		if failfast {
			logger.WarningTest(test.getMessage())
			continue
		}
		if err != nil {
			failed := test.failure(argoOptions)
			if failed {
				logger.FailureTest(test.getMessage())
				return err
			}
		}
		logger.SuccessTest(test.getMessage())
	}

	return nil
}

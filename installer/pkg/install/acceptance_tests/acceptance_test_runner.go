package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
)

type (
	acceptanceTest interface {
		Check(argoOptions *install.ArgoOptions) error
		GetMessage() string
	}

	IAcceptanceTestRunner interface {
		Verify(argoOptions *install.ArgoOptions) error
	}

	AcceptanceTestRunner struct {
	}
)

var tests []acceptanceTest
var runner IAcceptanceTestRunner

func New() IAcceptanceTestRunner {
	if runner == nil {
		tests = append(tests, &ApplicationAcceptanceTest{})
		tests = append(tests, &ProjectAcceptanceTest{})
		tests = append(tests, &ArgoCredentialsAcceptanceTest{})

		runner = AcceptanceTestRunner{}
	}
	return runner
}

func (runner AcceptanceTestRunner) Verify(argoOptions *install.ArgoOptions) error {
	logger.Info("\nTesting requirements")
	logger.Info("--------------------")
	defer logger.Info("--------------------\n")

	var err error

	for _, test := range tests {
		err = test.Check(argoOptions)
		if err != nil {
			logger.FailureTest(test.GetMessage())
			return err
		}
		logger.SuccessTest(test.GetMessage())
	}

	return nil
}

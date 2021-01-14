package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
)

type (
	acceptanceTest interface {
		Check(argoOptions *install.ArgoOptions) error
		GetMessage() string
	}

	IAcceptanceTestRunner interface {
		VerifyAgentSetup(argoOptions *install.ArgoOptions) error
		VerifyArgoSetup(argoOptions *install.ArgoOptions) error
	}

	AcceptanceTestRunner struct {
	}
)

var tests []acceptanceTest
var argoCredentialsTests []acceptanceTest

var runner IAcceptanceTestRunner

func New() IAcceptanceTestRunner {
	if runner == nil {

		argoCredentialsTests = append(argoCredentialsTests, &ArgoCredentialsAcceptanceTest{})

		argoCredentialsTests = append(argoCredentialsTests, &ProjectAcceptanceTest{
			argoApi: argo.GetInstance(),
		})
		argoCredentialsTests = append(argoCredentialsTests, &ApplicationAcceptanceTest{
			argoApi: argo.GetInstance(),
		})

		runner = AcceptanceTestRunner{}
	}
	return runner
}

func (runner AcceptanceTestRunner) verify(argoOptions *install.ArgoOptions, testsToExecute []acceptanceTest, title string) error {
	logger.Info("\n" + title)
	logger.Info("--------------------")
	defer logger.Info("--------------------\n")

	var err error

	for _, test := range testsToExecute {
		err = test.Check(argoOptions)
		if err != nil {
			logger.FailureTest(test.GetMessage())
			return err
		}
		logger.SuccessTest(test.GetMessage())
	}

	return nil
}

func (runner AcceptanceTestRunner) VerifyAgentSetup(argoOptions *install.ArgoOptions) error {
	return runner.verify(argoOptions, tests, "Testing requirements")
}

func (runner AcceptanceTestRunner) VerifyArgoSetup(argoOptions *install.ArgoOptions) error {
	return runner.verify(argoOptions, argoCredentialsTests, "Testing argocd credentials and permissions")
}

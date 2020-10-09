package questionnaire

import (
	"encoding/base64"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func AskAboutAgentPlaceOptions(installOptions *install.InstallCmdOptions) error {
	err, inCluster := prompt.Confirm("Is your Argo CD installation running on this cluster?")
	if err != nil {
		return err
	}

	installOptions.Kube.InCluster = inCluster

	if !inCluster {
		err = prompt.InputWithDefault(&installOptions.Kube.MasterUrl, "Enter Kubernetes URL where your Argo CD installation is running", "")
		if err != nil {
			return err
		}

		err = prompt.InputWithDefault(&installOptions.Kube.BearerToken, "Enter Kuberentes token where your Argo CD installation is running", "")
		if err != nil {
			return err
		}

		installOptions.Kube.BearerToken = base64.StdEncoding.EncodeToString([]byte(installOptions.Kube.BearerToken))
	}

	return nil
}

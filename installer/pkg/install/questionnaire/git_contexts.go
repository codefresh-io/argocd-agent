package questionnaire

import (
	"encoding/base64"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func AskAboutGitContext(installOptions *install.InstallCmdOptions) error {
	if installOptions.Git.Integration != "" { // Integration is passed
		err, context := holder.ApiHolder.GetGitContextByName(installOptions.Git.Integration)
		if err != nil {
			return err
		}
		installOptions.Git.Password = (*context).Spec.Data.Auth.Password
		return nil
	}

	err, contexts := holder.ApiHolder.GetGitContexts()
	if err != nil {
		return err
	}

	var values = make(map[string]string)
	var list []string
	for _, v := range *contexts {
		values[v.Metadata.Name] = v.Spec.Data.Auth.Password
		list = append(list, v.Metadata.Name)
	}

	err, selectedContext := prompt.Select(list, "Select Git context")
	installOptions.Git.Password = base64.StdEncoding.EncodeToString([]byte(values[selectedContext]))

	return err
}

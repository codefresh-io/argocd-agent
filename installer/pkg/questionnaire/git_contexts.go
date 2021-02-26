package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

func AskAboutGitContext(installOptions *install.InstallCmdOptions) error {
	if installOptions.Git.Integration != "" { // Integration is passed
		err, _ := holder.ApiHolder.GetGitContextByName(installOptions.Git.Integration)
		if err != nil {
			return err
		}
		return nil
	}

	err, contexts := holder.ApiHolder.GetGitContexts()
	if err != nil {
		return err
	}

	//var selectedContext string
	var values = make(map[string]string)
	var list []string
	for _, v := range *contexts {
		values[v.Metadata.Name] = v.Spec.Data.Auth.Password
		list = append(list, v.Metadata.Name)
	}

	if len(list) == 1 {
		installOptions.Git.Integration = list[0]
	} else {
		err, installOptions.Git.Integration = prompt.Select(list, "Select Git context (Please create a dedicated context for the agent to  avoid hitting the Github rate limits)")
	}

	return err
}
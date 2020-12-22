package git

import "github.com/codefresh-io/go-sdk/pkg/codefresh"

func GetAvailableContexts(cfContextsApi codefresh.IContextAPI) (*[]codefresh.ContextPayload, error) {
	var result = []codefresh.ContextPayload{}

	err, allContexts := cfContextsApi.GetGitContexts()
	if err != nil {
		return &result, err
	}
	for _, context := range *allContexts {
		//context.Spec.Data.Auth.SshPrivateKey != "" &&
		if context.Spec.Data.Auth.Type == "basic" && context.Spec.Type == "git.github" {
			result = append(result, context)
		}
	}
	return &result, nil
}

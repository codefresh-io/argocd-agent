package questionnaire

import (
	"encoding/base64"
	"encoding/json"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/util"
	"github.com/elliotchance/orderedmap"
)

// AskAboutSyncOptions ask about specific sync mode
func AskAboutSyncOptions(installOptions *entity.InstallCmdOptions) {
	var syncMode interface{}

	if installOptions.Codefresh.SyncMode != "" {
		syncMode = installOptions.Codefresh.SyncMode
	} else {
		syncModes := orderedmap.NewOrderedMap()
		syncModes.Set("Import all existing Argo applications to Codefresh", "SYNC")
		syncModes.Set("Select specific Argo applications to import", codefresh.SelectSync)
		syncModes.Set("Do not import anything from Argo to Codefresh", codefresh.None)

		_, autoSyncMode := prompt.NewPrompt().Select(util.ConvertIntToStringArray(syncModes.Keys()), "Select argocd sync behavior please")

		syncMode, _ = syncModes.Get(autoSyncMode)

		if syncMode == "SYNC" {
			_, autoSync := prompt.NewPrompt().Confirm("Enable auto-sync of applications, this will import all existing applications and update Codefresh in the future")
			if autoSync {
				syncMode = codefresh.ContinueSync
			} else {
				syncMode = codefresh.OneTimeSync
			}
		}
	}

	if syncMode == codefresh.SelectSync {
		applicationsForSync := installOptions.Codefresh.ApplicationsForSyncArr

		if len(applicationsForSync) == 0 {

			argoToken := installOptions.Argo.Token

			if installOptions.Argo.Username != "" {
				argoToken, _ = argo.GetUnauthorizedApiInstance().GetToken(installOptions.Argo.Username, installOptions.Argo.Password, installOptions.Argo.Host)
			}

			applications, _ := argo.GetUnauthorizedApiInstance().GetApplications(argoToken, installOptions.Argo.Host)

			applicationNames := make([]string, 0)

			for _, prj := range applications {
				applicationNames = append(applicationNames, prj.Metadata.Name)
			}

			_, applicationsForSync = prompt.NewPrompt().Multiselect(applicationNames, "Please select application for sync")
		}

		applicationsAsJson, _ := json.Marshal(applicationsForSync)

		installOptions.Codefresh.ApplicationsForSync = base64.StdEncoding.EncodeToString(applicationsAsJson)
	}

	installOptions.Codefresh.SyncMode = syncMode.(string)
}

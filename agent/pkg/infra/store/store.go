package store

var (
	store *Values
)

type Environment struct {
	Name string
}

type (
	Values struct {
		Agent struct {
			Version string
		}
		Git struct {
			Token       string
			Integration string
		}
		Argo struct {
			Token string
			Host  string
		}
		Codefresh struct {
			Host                string
			Token               string
			Integration         string
			SyncMode            string
			ApplicationsForSync []string
		}
		Heartbeat struct {
			Error string
		}
		Environments []Environment
	}
)

func GetStore() *Values {
	if store == nil {
		store = &Values{}
		return store
	}
	return store
}

func SetArgo(token string, host string) *Values {
	values := GetStore()
	values.Argo.Token = token
	values.Argo.Host = host
	return values
}

func SetCodefresh(host string, token string, integration string) *Values {
	values := GetStore()
	values.Codefresh.Token = token
	values.Codefresh.Host = host
	values.Codefresh.Integration = integration
	return values
}

func SetSyncOptions(syncMode string, applicationsToSync []string) *Values {
	values := GetStore()
	values.Codefresh.SyncMode = syncMode
	values.Codefresh.ApplicationsForSync = applicationsToSync
	return values
}

func SetHeartbeatError(error string) *Values {
	values := GetStore()
	values.Heartbeat.Error = error
	return values
}

func SetEnvironments(environments []Environment) *Values {
	values := GetStore()
	values.Environments = environments
	return values
}
func SetGit(Token string) *Values {
	values := GetStore()
	values.Git.Token = Token
	return values
}

func SetAgent(Version string) *Values {
	values := GetStore()
	values.Agent.Version = Version
	return values
}

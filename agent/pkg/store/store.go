package store

var (
	store *Values
)

type Environment struct {
	Name string
}

type (
	Values struct {
		Git struct {
			Token string
		}
		Argo struct {
			Token string
			Host  string
		}
		Codefresh struct {
			Host        string
			Token       string
			Integration string
			AutoSync    bool
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

func SetCodefresh(host string, token string, integration string, autoSync bool) *Values {
	values := GetStore()
	values.Codefresh.Token = token
	values.Codefresh.Host = host
	values.Codefresh.Integration = integration
	values.Codefresh.AutoSync = autoSync
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

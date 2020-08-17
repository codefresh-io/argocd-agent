package store

var (
	store *Values
)

type Environment struct {
	Name string
}

type (
	Values struct {
		Argo struct {
			Token     string
			Host      string
			Namespace string
		}
		Codefresh struct {
			Host        string
			Token       string
			Integration string
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

func SetArgo(token string, host string, namespace string) *Values {
	values := GetStore()
	values.Argo.Token = token
	values.Argo.Host = host
	values.Argo.Namespace = namespace
	return values
}

func SetCodefresh(host string, token string, integration string) *Values {
	values := GetStore()
	values.Codefresh.Token = token
	values.Codefresh.Host = host
	values.Codefresh.Integration = integration
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

package store

var (
	store *Values
)

type (
	Values struct {
		Argo struct {
			Token string
			Host  string
		}
		Codefresh struct {
			Host        string
			Token       string
			Integration string
		}
		Heartbeat struct {
			Error string
		}
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

func SetHeartbeatError(error string) *Values {
	values := GetStore()
	values.Heartbeat.Error = error
	return values
}

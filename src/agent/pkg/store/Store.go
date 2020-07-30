package store

var (
	store *Values
)

type (
	Values struct {
		Argo struct {
			Token string
		}
		Codefresh struct {
			Host  string
			Token string
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

func SetArgoToken(token string) *Values {
	values := GetStore()
	values.Argo.Token = token
	return values
}

func SetCodefresh(host string, token string) *Values {
	values := GetStore()
	values.Codefresh.Token = token
	values.Codefresh.Host = host
	return values
}

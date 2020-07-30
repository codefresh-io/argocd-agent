package store

var (
	store *Values
)

type (
	Values struct {
		Token string
	}
)

func GetStore() *Values {
	if store == nil {
		store = &Values{}
		return store
	}
	return store
}

func SetToken(token string) *Values {
	values := GetStore()
	values.Token = token
	return values
}

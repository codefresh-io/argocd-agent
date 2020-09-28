package git

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Annotation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Gitops struct {
	Comitters []User       `json:"comitters"`
	Prs       []Annotation `json:"prs"`
	Issues    []Annotation `json:"issues"`
}

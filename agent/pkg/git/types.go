package git

type User struct {
	Name   string
	Avatar string
}

type Annotation struct {
	Key   string
	Value string
}

type Gitops struct {
	Committers []User       `json:"committers"`
	Prs        []Annotation `json:"prs"`
	Issues     []Annotation `json:"issues"`
}

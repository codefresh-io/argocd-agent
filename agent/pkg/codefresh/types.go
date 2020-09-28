package codefresh

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/guregu/null"
)

type MongoCFEnvWrapper struct {
	Docs []CFEnvironment `json:"docs"`
}

type CFEnvironment struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Type        string `json:"type"`
		Application string `json:"application"`
	} `json:"spec"`
}

type GitInfo struct {
	Committers []*github.User        `json:"committers"`
	Prs        []*github.PullRequest `json:"prs"`
	Issues     []*github.Issue       `json:"issues"`
}

type Environment struct {
	GitInfo      GitInfo               `json:"gitInfo"`
	FinishedAt   string                `json:"finishedAt"`
	HealthStatus string                `json:"healthStatus"`
	SyncStatus   string                `json:"status"`
	HistoryId    int64                 `json:"historyId"`
	SyncRevision string                `json:"revision"`
	Name         string                `json:"name"`
	Activities   []EnvironmentActivity `json:"activities"`
	Resources    interface{}           `json:"resources"`
	RepoUrl      string                `json:"repoUrl"`
}

type EnvironmentActivity struct {
	Name         string   `json:"name"`
	TargetImages []string `json:"targetImages"`
	Status       string   `json:"status"`
	LiveImages   []string `json:"liveImages"`
}

type CodefreshError struct {
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Name    string      `json:"name"`
	Message string      `json:"message"`
	Context interface{} `json:"context"`
	URL     string
}

func (e *CodefreshError) Error() string {
	return fmt.Sprintf("Request failed to %s, %s - %s", e.URL, e.Code, e.Message)
}

type AgentApplication struct {
	Name      string `json:"name"`
	UID       string `json:"uid"`
	Project   string `json:"project"`
	Namespace string `json:"namespace"`
	Server    string `json:"server"`
}

type AgentProject struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

type AgentState struct {
	Kind  string      `json:"type"`
	Items interface{} `json:"items"`
}

type IntegrationPayloadData struct {
	Name     string      `json:"name"`
	Url      string      `json:"url"`
	Username null.String `json:"username"`
	Password null.String `json:"password"`
	Token    null.String `json:"token"`
}

type IntegrationPayload struct {
	Type string                 `json:"type"`
	Data IntegrationPayloadData `json:"data"`
}

type EnvironmentMetadata struct {
	Name string `json:"name"`
}

type EnvironmentSpec struct {
	Type        string `json:"type"`
	Context     string `json:"context"`
	Project     string `json:"project"`
	Application string `json:"application"`
}

type EnvironmentPayload struct {
	Version  string              `json:"version"`
	Metadata EnvironmentMetadata `json:"metadata"`
	Spec     EnvironmentSpec     `json:"spec"`
}

type Heartbeat struct {
	Error string `json:"error"`
}

type requestOptions struct {
	path   string
	method string
	body   interface{}
	qs     map[string]string
}

type ContextPayload struct {
	Spec struct {
		Type string `json:"type"`
		Data struct {
			Auth struct {
				Password      string `json:"password"`
				ApiHost       string `json:"apiHost"`
				ApiPathPrefix string `json:"apiPathPrefix"`
			} `json:"auth"`
		} `json:"data"`
	} `json:"spec"`
}

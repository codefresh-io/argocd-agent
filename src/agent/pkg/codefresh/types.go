package codefresh

type Environment struct {
	FinishedAt   string                `json:"finishedAt"`
	HealthStatus string                `json:"healthStatus"`
	SyncStatus   string                `json:"status"`
	SyncRevision string                `json:"revision"`
	Name         string                `json:"name"`
	Activities   []EnvironmentActivity `json:"activities"`
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
}

type AgentApplication struct {
	Name    string `json:"name"`
	UID     string `json:"uid"`
	Project string `json:"project"`
}

type AgentProject struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

type AgentState struct {
	Kind  string      `json:"type"`
	Items interface{} `json:"items"`
}

type requestOptions struct {
	path   string
	method string
	body   interface{}
}

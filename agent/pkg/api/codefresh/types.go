package codefresh

import (
	"fmt"
)

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

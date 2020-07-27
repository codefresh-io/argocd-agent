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

package events

// EventHandler is interface for handle different type of events, like new application added, rollout happens so on
type EventHandler interface {
	Handle(payload interface{}) error
}

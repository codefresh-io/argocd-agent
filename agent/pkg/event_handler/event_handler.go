package event_handler

type EventHandler interface {
	Handle(payload interface{}) error
}

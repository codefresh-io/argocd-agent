package cf_event_sender

import (
	"testing"
)

func TestNew(t *testing.T) {
	client := New(EVENT_UNINSTALL)
	if client.eventName != EVENT_UNINSTALL {
		t.Errorf("'TestNew' failed, unexpected event name after init, expected '%v', got '%v'", EVENT_UNINSTALL, client.eventName)
	}

	client = New(EVENT_INSTALL)
	if client.eventName != EVENT_UNINSTALL {
		t.Errorf("'TestNew' failed, must return existing state, expected '%v', got '%v'", EVENT_UNINSTALL, client.eventName)
	}
}

package cfeventsender

import (
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestNew(t *testing.T) {
	client := New(EVENT_AGENT_UNINSTALL)
	if client.eventName != EVENT_AGENT_UNINSTALL {
		t.Errorf("'TestNew' failed, unexpected event name after init, expected '%v', got '%v'", EVENT_AGENT_UNINSTALL, client.eventName)
	}

	client = New(EVENT_AGENT_INSTALL)
	if client.eventName != EVENT_AGENT_INSTALL {
		t.Errorf("'TestNew' failed, must return existing state, expected '%v', got '%v'", EVENT_AGENT_INSTALL, client.eventName)
	}
}

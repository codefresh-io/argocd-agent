package cf_event_sender

import "github.com/codefresh-io/argocd-listener/installer/pkg/holder"

const (
	STATUS_SUCCESS  = "Success"
	STATUS_FAILED   = "Failed"
	EVENT_UNINSTALL = "agent.uninstalled"
	EVENT_INSTALL   = "agent.installed"
)

type CfEventSender struct {
	eventName string
}

var cfEventSender *CfEventSender

func New(eventName string) *CfEventSender {
	if cfEventSender == nil {
		cfEventSender = &CfEventSender{eventName}
	}
	return cfEventSender
}

func (cfEventSender *CfEventSender) Send(status string, reason string) {
	props := make(map[string]string)
	props["status"] = status
	props["reason"] = reason
	_ = holder.ApiHolder.SendEvent(cfEventSender.eventName, props)
}

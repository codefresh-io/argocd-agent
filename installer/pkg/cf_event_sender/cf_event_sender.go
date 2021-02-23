package cf_event_sender

import "github.com/codefresh-io/argocd-listener/installer/pkg/holder"

const (
	STATUS_SUCCESS = "Success"
	STATUS_FAILED  = "Failed"

	EVENT_AGENT_INSTALL   = "agent.installed"
	EVENT_AGENT_UNINSTALL = "agent.uninstalled"

	EVENT_CONTROLLER_INSTALL   = "controller.installed"
	EVENT_CONTROLLER_UNINSTALL = "controller.uninstalled"
)

type CfEventSender struct {
	eventName string
}

func New(eventName string) *CfEventSender {
	return &CfEventSender{eventName}
}

func (cfEventSender *CfEventSender) Success(reason string) {
	props := make(map[string]string)
	props["status"] = STATUS_SUCCESS
	props["reason"] = reason
	_ = holder.ApiHolder.SendEvent(cfEventSender.eventName, props)
}

func (cfEventSender *CfEventSender) Fail(reason string) {
	props := make(map[string]string)
	props["status"] = STATUS_FAILED
	props["reason"] = reason
	_ = holder.ApiHolder.SendEvent(cfEventSender.eventName, props)
}

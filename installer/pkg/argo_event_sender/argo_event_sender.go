package argo_event_sender

import "github.com/codefresh-io/argocd-listener/installer/pkg/holder"

const (
	STATUS_SUCCESS  = "Success"
	STATUS_FAILED   = "Failed"
	EVENT_UNINSTALL = "agent.uninstalled"
	EVENT_INSTALL   = "agent.installed"
)

type ArgoEventSender struct {
	eventName string
}

var argoEventSender *ArgoEventSender

func New(eventName string) *ArgoEventSender {
	if argoEventSender == nil {
		argoEventSender = &ArgoEventSender{eventName}
	}
	return argoEventSender
}

func (argoEventSender *ArgoEventSender) Send(status string, reason string) {
	props := make(map[string]string)
	props["status"] = status
	props["reason"] = reason
	_ = holder.ApiHolder.SendEvent(argoEventSender.eventName, props)
}

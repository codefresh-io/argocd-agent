package handler

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	handler "github.com/codefresh-io/argocd-listener/installer/pkg/install/installer/strategy"
)

type (
	installerStrategy interface {
		Install(installCmdOptions install.InstallCmdOptions) (error, string)
	}

	IInstaller interface {
		Install(installCmdOptions install.InstallCmdOptions) (error, string)
	}

	Installer struct {
	}
)

var installer IInstaller

func New() IInstaller {
	if installer == nil {
		installer = Installer{}
	}
	return installer
}

func (i Installer) Install(installCmdOptions install.InstallCmdOptions) (error, string) {
	var strategy installerStrategy
	if installCmdOptions.Agent.Install {
		strategy = &handler.RealArgoCdAgentInstaller{}
	} else {
		strategy = &handler.DryRunArgoCdAgentInstaller{}
	}

	return strategy.Install(installCmdOptions)
}

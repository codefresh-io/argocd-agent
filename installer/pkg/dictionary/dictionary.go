package dictionary

const (
	NoApplication                = "Found no applications in your argocd installation"
	CouldNotAccessToArgocdServer = "could not access your argocd server %v"

	CheckArgoApplicationsAccessability = "checking argocd applications accessibility..."
	CheckArgoServerAccessability       = "checking argocd server accessibility..."
	CheckArgoCredentials               = "checking argocd credentials..."

	// Applications acceptance test
	StopInstallation     = "Stop installation"
	ContinueInstallation = "Continue installation (you can add later applications manually)"
	SetupDemoApplication = "Setup demo application in ArgoCD (viewable also in Codefresh)"

	// Argo server acceptance test
	ContinueInstallationBehindFirewall = "Continue installation (all future acceptance tests will be skipped)"
)

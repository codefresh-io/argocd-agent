package dictionary

const (
	NoApplication                = "Found no applications in your ArgoCD installation"
	CouldNotAccessToArgocdServer = "could not access your ArgoCD server %v"

	CheckArgoApplicationsAccessability = "checking ArgoCD applications accessibility..."
	CheckArgoServerAccessability       = "checking ArgoCD server accessibility..."
	CheckArgoCredentials               = "checking ArgoCD credentials..."

	// Applications acceptance test
	StopInstallation     = "Stop installation"
	ContinueInstallation = "Continue installation (you can add later applications manually)"
	SetupDemoApplication = "Setup demo application in ArgoCD (viewable also in Codefresh)"

	// Argo server acceptance test
	ContinueInstallationBehindFirewall = "Continue installation (all future acceptance tests will be skipped)"
)

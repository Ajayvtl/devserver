package capabilities

// Capability names what a provider can do.
type Capability string

const (
	CapabilityInstall    Capability = "install"
	CapabilityDeploy     Capability = "deploy"
	CapabilityBuild      Capability = "build"
	CapabilityAnalyze    Capability = "analyze"
	CapabilityBackup     Capability = "backup"
	CapabilityRestore    Capability = "restore"
	CapabilityDatabase   Capability = "database"
	CapabilityAI         Capability = "ai"
	CapabilityTerminal   Capability = "terminal"
	CapabilityLogs       Capability = "logs"
	CapabilityMonitoring Capability = "monitoring"
	CapabilitySecrets    Capability = "secrets"
	CapabilityFilesystem Capability = "filesystem"
)

func (c Capability) String() string {
	return string(c)
}

package contracts

// Capabilities is a typed capability model instead of boolean flags.
type Capabilities struct {
	Shell                *ShellCapability
	ServiceControl       *ServiceControlCapability
	Docker               *DockerCapability
	Kubernetes           *KubernetesCapability
	Filesystem           *FilesystemCapability
	PortForward          *PortForwardCapability
	Process              *ProcessCapability
	EnvironmentVariables *EnvVarCapability
	Git                  *GitCapability
	PackageManager       *PackageManagerCapability
}

type ShellCapability struct {
	Supported bool
	Type      string // bash, zsh, powershell
}

type ServiceControlCapability struct {
	Supported bool
	Type      string // systemd, launchd, windows_scm
}

type DockerCapability struct {
	Supported bool
	Version   string
	Features  []string // compose, swarm, buildkit
}

type KubernetesCapability struct {
	Supported bool
	Version   string
}

type FilesystemCapability struct {
	Supported bool
}

type PortForwardCapability struct {
	Supported bool
}

type ProcessCapability struct {
	Supported bool
}

type EnvVarCapability struct {
	Supported bool
}

type GitCapability struct {
	Supported bool
	Version   string
}

type PackageManagerCapability struct {
	Supported bool
	Type      string // apt, yum, brew, npm
}

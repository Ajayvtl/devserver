package capabilities

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ajayvtl/devserver/internal/registry"
)

// ModuleProvider adapts an existing registry.Module to the capability system.
type ModuleProvider struct {
	module registry.Module
}

func NewModuleProvider(module registry.Module) *ModuleProvider {
	return &ModuleProvider{module: module}
}

func (p *ModuleProvider) Name() string {
	if p == nil || p.module == nil {
		return ""
	}
	return p.module.Name()
}

func (p *ModuleProvider) Capabilities() []Capability {
	if p == nil || p.module == nil {
		return nil
	}
	return capabilitiesForModule(p.module.Name())
}

func (p *ModuleProvider) Metadata(cap Capability) Metadata {
	return Metadata{
		Provider:    p.Name(),
		Capability:  cap,
		Title:       titleForCapability(cap),
		Summary:     p.module.Description(),
		Category:    categoryForCapability(cap),
		Priority:    priorityForModule(p.Name(), cap),
		Tags:        []string{p.Name(), string(cap)},
		Permissions: permissionsList(permissionsForCapability(cap)),
	}
}

func (p *ModuleProvider) Permissions(cap Capability) PermissionSet {
	return permissionsForCapability(cap)
}

func (p *ModuleProvider) Execute(ctx context.Context, cap Capability, payload any) error {
	if p == nil || p.module == nil {
		return fmt.Errorf("module provider is unavailable")
	}

	switch cap {
	case CapabilityInstall:
		return p.module.Install(ctx)
	case CapabilityDeploy:
		return p.module.Configure(ctx)
	case CapabilityBuild:
		return p.module.Validate(ctx)
	case CapabilityAnalyze:
		return p.module.Check(ctx)
	case CapabilityBackup:
		return p.module.Upgrade(ctx)
	case CapabilityRestore:
		return p.module.Rollback(ctx)
	case CapabilityDatabase:
		return p.module.Check(ctx)
	case CapabilityAI:
		return p.module.Validate(ctx)
	case CapabilityTerminal:
		return p.module.Configure(ctx)
	case CapabilityLogs:
		return p.module.Check(ctx)
	case CapabilityMonitoring:
		return p.module.Validate(ctx)
	case CapabilitySecrets:
		return p.module.Check(ctx)
	case CapabilityFilesystem:
		return p.module.Validate(ctx)
	default:
		return fmt.Errorf("provider %q does not implement capability %q", p.module.Name(), cap)
	}
}

func (p *ModuleProvider) Health(ctx context.Context) error {
	if p == nil || p.module == nil {
		return fmt.Errorf("module provider is unavailable")
	}
	return p.module.Check(ctx)
}

func capabilitiesForModule(name string) []Capability {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "nginx":
		return []Capability{CapabilityDeploy, CapabilityLogs, CapabilityMonitoring, CapabilityInstall, CapabilityAnalyze}
	case "node":
		return []Capability{CapabilityBuild, CapabilityAnalyze, CapabilityTerminal, CapabilityLogs, CapabilityInstall}
	case "postgres", "mysql", "redis":
		return []Capability{CapabilityDatabase, CapabilityBackup, CapabilityRestore, CapabilityLogs, CapabilityMonitoring, CapabilityInstall}
	case "php", "python":
		return []Capability{CapabilityBuild, CapabilityAnalyze, CapabilityTerminal, CapabilityLogs, CapabilityInstall}
	default:
		return []Capability{CapabilityAnalyze, CapabilityLogs, CapabilityInstall}
	}
}

func titleForCapability(cap Capability) string {
	switch cap {
	case CapabilityInstall:
		return "Install"
	case CapabilityDeploy:
		return "Deploy"
	case CapabilityBuild:
		return "Build"
	case CapabilityAnalyze:
		return "Analyze"
	case CapabilityBackup:
		return "Backup"
	case CapabilityRestore:
		return "Restore"
	case CapabilityDatabase:
		return "Database"
	case CapabilityAI:
		return "AI"
	case CapabilityTerminal:
		return "Terminal"
	case CapabilityLogs:
		return "Logs"
	case CapabilityMonitoring:
		return "Monitoring"
	case CapabilitySecrets:
		return "Secrets"
	case CapabilityFilesystem:
		return "Filesystem"
	default:
		if cap == "" {
			return "Capability"
		}
		return strings.ToUpper(string(cap[:1])) + string(cap[1:])
	}
}

func categoryForCapability(cap Capability) string {
	switch cap {
	case CapabilityInstall, CapabilityDeploy, CapabilityBuild:
		return "operations"
	case CapabilityAnalyze, CapabilityLogs, CapabilityMonitoring:
		return "observability"
	case CapabilityBackup, CapabilityRestore, CapabilityDatabase:
		return "data"
	case CapabilityAI:
		return "intelligence"
	case CapabilityTerminal, CapabilityFilesystem, CapabilitySecrets:
		return "platform"
	default:
		return "general"
	}
}

func permissionsForCapability(cap Capability) PermissionSet {
	switch cap {
	case CapabilityDeploy, CapabilityTerminal:
		return NewPermissionSet(PermissionExecute, PermissionNetwork)
	case CapabilityBuild, CapabilityAnalyze, CapabilityLogs, CapabilityMonitoring:
		return NewPermissionSet(PermissionRead)
	case CapabilityBackup, CapabilityRestore, CapabilityDatabase:
		return NewPermissionSet(PermissionRead, PermissionWrite)
	case CapabilityAI:
		return NewPermissionSet(PermissionNetwork, PermissionSecrets)
	case CapabilitySecrets:
		return NewPermissionSet(PermissionSecrets)
	case CapabilityFilesystem:
		return NewPermissionSet(PermissionFilesystem)
	default:
		return NewPermissionSet(PermissionRead)
	}
}

func permissionsList(set PermissionSet) []Permission {
	out := make([]Permission, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	return out
}

func priorityForModule(name string, cap Capability) int {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "postgres":
		return 90
	case "mysql":
		return 80
	case "redis":
		return 70
	case "nginx":
		return 85
	case "node":
		return 75
	case "php":
		return 60
	case "python":
		return 65
	}

	switch cap {
	case CapabilityDeploy:
		return 80
	case CapabilityDatabase:
		return 75
	case CapabilityBuild:
		return 70
	case CapabilityAnalyze:
		return 60
	default:
		return 50
	}
}

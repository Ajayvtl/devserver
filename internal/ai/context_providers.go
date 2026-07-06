package ai

import (
	"fmt"
	"strings"

	"github.com/Ajayvtl/devserver/internal/core"
)

// ContextProvider is responsible for extracting a specific piece of context
// and formatting it into a string block.
type ContextProvider interface {
	Name() string
	Provide(ws *core.WorkspaceContext, editor *EditorState) string
}

// EditorContextProvider extracts active file and selection data.
type EditorContextProvider struct{}

func (p *EditorContextProvider) Name() string { return "Editor" }
func (p *EditorContextProvider) Provide(ws *core.WorkspaceContext, editor *EditorState) string {
	if editor == nil {
		return ""
	}
	var b strings.Builder
	if editor.ActiveFile != "" {
		b.WriteString(fmt.Sprintf("Active File: %s\n", editor.ActiveFile))
	}
	if editor.Language != "" {
		b.WriteString(fmt.Sprintf("Language: %s\n", editor.Language))
	}
	if editor.SelectedText != "" {
		b.WriteString(fmt.Sprintf("Selected Text:\n```\n%s\n```\n", editor.SelectedText))
	}
	return b.String()
}

// DiagnosticsProvider extracts active file diagnostics and workspace health.
type DiagnosticsProvider struct{}

func (p *DiagnosticsProvider) Name() string { return "Diagnostics" }
func (p *DiagnosticsProvider) Provide(ws *core.WorkspaceContext, editor *EditorState) string {
	var b strings.Builder
	if editor != nil && len(editor.Diagnostics) > 0 {
		b.WriteString("Active File Diagnostics:\n")
		// Budget optimization: limit to top 10
		limit := len(editor.Diagnostics)
		if limit > 10 {
			limit = 10
		}
		for i := 0; i < limit; i++ {
			d := editor.Diagnostics[i]
			b.WriteString(fmt.Sprintf("- [%s] %s:%d: %s\n", d.Severity, d.File, d.Line, d.Message))
		}
		if len(editor.Diagnostics) > 10 {
			b.WriteString(fmt.Sprintf("... (and %d more diagnostics hidden)\n", len(editor.Diagnostics)-10))
		}
	}

	if ws != nil && (ws.Health.Build != "" || ws.Health.Tests != "" || ws.Health.Lint != "") {
		b.WriteString("Workspace Health:\n")
		if ws.Health.Build != "" && ws.Health.Build != "unknown" {
			b.WriteString(fmt.Sprintf("- Build: %s\n", ws.Health.Build))
		}
		if ws.Health.Tests != "" && ws.Health.Tests != "unknown" {
			b.WriteString(fmt.Sprintf("- Tests: %s\n", ws.Health.Tests))
		}
		if ws.Health.Lint != "" && ws.Health.Lint != "unknown" {
			b.WriteString(fmt.Sprintf("- Lint: %s\n", ws.Health.Lint))
		}
	}
	return b.String()
}

// WorkspaceProvider extracts core workspace metadata and topology.
type WorkspaceProvider struct{}

func (p *WorkspaceProvider) Name() string { return "Workspace" }
func (p *WorkspaceProvider) Provide(ws *core.WorkspaceContext, editor *EditorState) string {
	if ws == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Workspace: %s (%s)\n", ws.Workspace.Name, ws.Workspace.Kind))

	if len(ws.Services) > 0 {
		b.WriteString("Active Services: ")
		var srvNames []string
		for _, s := range ws.Services {
			if s.State.Status == "running" {
				srvNames = append(srvNames, s.Name)
			}
		}
		b.WriteString(strings.Join(srvNames, ", "))
		b.WriteString("\n")
	}

	if len(ws.Knowledge.Symbols) > 0 {
		b.WriteString(fmt.Sprintf("Known Symbols: %d\n", len(ws.Knowledge.Symbols)))
	}
	return b.String()
}

// GitProvider extracts repository status.
type GitProvider struct{}

func (p *GitProvider) Name() string { return "Git" }
func (p *GitProvider) Provide(ws *core.WorkspaceContext, editor *EditorState) string {
	if ws == nil || ws.Git.Branch == "" {
		return ""
	}
	return fmt.Sprintf("Git Branch: %s (Commit: %s)\n", ws.Git.Branch, ws.Git.Commit)
}

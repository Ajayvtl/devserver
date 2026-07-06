package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ajayvtl/devserver/internal/core"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/rs/zerolog"
)

// Diagnostic represents a compiler error, linter warning, or LSP diagnostic.
type Diagnostic struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// EditorState represents the context provided by the UI Editor.
type EditorState struct {
	ActiveFile   string       `json:"activeFile,omitempty"`
	Language     string       `json:"language,omitempty"`
	Cursor       float64      `json:"cursor,omitempty"`
	SelectedText string       `json:"selectedText,omitempty"`
	Diagnostics  []Diagnostic `json:"diagnostics,omitempty"`
}

// AssembledContext represents the compiled state of the workspace and editor for AI inference.
type AssembledContext struct {
	WorkspaceContext *core.WorkspaceContext
	EditorState      *EditorState
	Summary          string
}

// ContextAssembler is responsible for gathering all static and dynamic state for the AI.
type ContextAssembler struct {
	log       zerolog.Logger
	workspace *core.WorkspaceProvider
}

// NewContextAssembler creates a new ContextAssembler.
func NewContextAssembler(logger zerolog.Logger, workspace *core.WorkspaceProvider) *ContextAssembler {
	return &ContextAssembler{
		log:       logger.With().Str("component", "ContextAssembler").Logger(),
		workspace: workspace,
	}
}

// Assemble pulls together the full workspace context for the given workspace ID.
func (c *ContextAssembler) Assemble(ctx context.Context, workspaceID common.WorkspaceID, editorState *EditorState) (*AssembledContext, error) {
	wsContext, err := c.workspace.Load(string(workspaceID))
	if err != nil {
		return nil, fmt.Errorf("failed to load workspace context: %w", err)
	}

	if editorState == nil {
		editorState = &EditorState{}
	}

	summary := c.generateSummary(wsContext, editorState)

	return &AssembledContext{
		WorkspaceContext: wsContext,
		EditorState:      editorState,
		Summary:          summary,
	}, nil
}

func (c *ContextAssembler) generateSummary(ws *core.WorkspaceContext, editor *EditorState) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Workspace: %s\n", ws.Workspace.Name))

	// Inject Editor Context if present
	if editor.ActiveFile != "" {
		builder.WriteString(fmt.Sprintf("Active File: %s\n", editor.ActiveFile))
	}
	if editor.Language != "" {
		builder.WriteString(fmt.Sprintf("Language: %s\n", editor.Language))
	}
	if editor.SelectedText != "" {
		builder.WriteString(fmt.Sprintf("Selected Text:\n```\n%s\n```\n", editor.SelectedText))
	}

	if len(editor.Diagnostics) > 0 {
		builder.WriteString("Diagnostics:\n")
		for _, d := range editor.Diagnostics {
			builder.WriteString(fmt.Sprintf("- [%s] %s:%d: %s\n", d.Severity, d.File, d.Line, d.Message))
		}
	}

	if ws.Git.Branch != "" {
		builder.WriteString(fmt.Sprintf("Git Branch: %s\n", ws.Git.Branch))
	}

	if len(ws.Services) > 0 {
		builder.WriteString("Active Services: ")
		var srvNames []string
		for _, s := range ws.Services {
			if s.State.Status == "running" {
				srvNames = append(srvNames, s.Name)
			}
		}
		builder.WriteString(strings.Join(srvNames, ", "))
		builder.WriteString("\n")
	}

	if len(ws.Knowledge.Symbols) > 0 {
		builder.WriteString(fmt.Sprintf("Known Symbols: %d\n", len(ws.Knowledge.Symbols)))
	}

	if ws.Health.Build != "" || ws.Health.Tests != "" || ws.Health.Lint != "" {
		builder.WriteString("Workspace Health (Diagnostics):\n")
		if ws.Health.Build != "" && ws.Health.Build != "unknown" {
			builder.WriteString(fmt.Sprintf("- Build: %s\n", ws.Health.Build))
		}
		if ws.Health.Tests != "" && ws.Health.Tests != "unknown" {
			builder.WriteString(fmt.Sprintf("- Tests: %s\n", ws.Health.Tests))
		}
		if ws.Health.Lint != "" && ws.Health.Lint != "unknown" {
			builder.WriteString(fmt.Sprintf("- Lint: %s\n", ws.Health.Lint))
		}
	}

	return builder.String()
}

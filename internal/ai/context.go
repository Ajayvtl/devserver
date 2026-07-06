package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ajayvtl/devserver/internal/core"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/rs/zerolog"
)

// EditorState represents the context provided by the UI Editor.
type EditorState struct {
	ActiveFile   string  `json:"activeFile,omitempty"`
	Language     string  `json:"language,omitempty"`
	Cursor       float64 `json:"cursor,omitempty"`
	SelectedText string  `json:"selectedText,omitempty"`
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

	builder.WriteString(fmt.Sprintf("Known Symbols: %d\n", len(ws.Knowledge.Symbols)))

	return builder.String()
}

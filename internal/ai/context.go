package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ajayvtl/devserver/internal/core"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/rs/zerolog"
)

// AssembledContext represents the compiled state of the workspace and editor for AI inference.
type AssembledContext struct {
	WorkspaceContext *core.WorkspaceContext
	EditorState      map[string]any // This will be expanded when wired to SessionManager
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
func (c *ContextAssembler) Assemble(ctx context.Context, workspaceID common.WorkspaceID) (*AssembledContext, error) {
	wsContext, err := c.workspace.Load(string(workspaceID))
	if err != nil {
		return nil, fmt.Errorf("failed to load workspace context: %w", err)
	}

	// In the future, we will also fetch the Editor Session (open tabs, cursor position)
	// from the editor.SessionManager via the EventBridge or direct dependency.

	summary := c.generateSummary(wsContext)

	return &AssembledContext{
		WorkspaceContext: wsContext,
		EditorState:      make(map[string]any),
		Summary:          summary,
	}, nil
}

func (c *ContextAssembler) generateSummary(ws *core.WorkspaceContext) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Workspace: %s\n", ws.Workspace.Name))

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

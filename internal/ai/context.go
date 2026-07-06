package ai

import (
	"context"
	"fmt"

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
	budget    *ContextBudgetManager
}

// NewContextAssembler creates a new ContextAssembler.
func NewContextAssembler(logger zerolog.Logger, workspace *core.WorkspaceProvider) *ContextAssembler {
	return &ContextAssembler{
		log:       logger.With().Str("component", "ContextAssembler").Logger(),
		workspace: workspace,
		budget:    NewContextBudgetManager(4000), // e.g. 4000 characters limit
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

	summary := c.budget.Build(wsContext, editorState)

	return &AssembledContext{
		WorkspaceContext: wsContext,
		EditorState:      editorState,
		Summary:          summary,
	}, nil
}

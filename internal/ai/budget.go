package ai

import (
	"strings"

	"github.com/Ajayvtl/devserver/internal/core"
)

// ContextBudgetManager orchestrates a pipeline of ContextProviders,
// assembling the final context block while enforcing token/character bounds.
type ContextBudgetManager struct {
	providers []ContextProvider
	maxChars  int
}

// NewContextBudgetManager creates a budget manager with the given providers.
func NewContextBudgetManager(limitChars int) *ContextBudgetManager {
	return &ContextBudgetManager{
		maxChars: limitChars,
		providers: []ContextProvider{
			&EditorContextProvider{},
			&DiagnosticsProvider{},
			&WorkspaceProvider{},
			&GitProvider{},
		},
	}
}

// Build compiles the outputs from all context providers, stopping early
// if the budget (character limit) is reached.
func (m *ContextBudgetManager) Build(ws *core.WorkspaceContext, editor *EditorState) string {
	var builder strings.Builder
	remaining := m.maxChars

	for _, p := range m.providers {
		block := p.Provide(ws, editor)
		if block == "" {
			continue
		}

		// Ensure we add spacing between blocks
		if builder.Len() > 0 {
			block = "\n" + block
		}

		if len(block) > remaining {
			// Trim to budget and indicate truncation
			builder.WriteString(block[:remaining])
			builder.WriteString("\n...[Context truncated due to budget]...\n")
			break
		}

		builder.WriteString(block)
		remaining -= len(block)
		if remaining <= 0 {
			break
		}
	}

	return builder.String()
}

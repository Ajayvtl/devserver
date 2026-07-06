package ai

import (
	"context"
	"fmt"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/tasks"
)

// TaskRunner implements the tasks.Runner interface to handle AI intents (generate, refactor, explain)
// triggered via the Command Bus.
type TaskRunner struct {
	AIRuntime *Runtime
}

// Execute processes the task payload and dispatches it to the AI Runtime.
func (r *TaskRunner) Execute(ctx context.Context, task *tasks.Task, runtime *tasks.Runtime) error {
	if r.AIRuntime == nil {
		return fmt.Errorf("AI Runtime is not initialized")
	}

	prompt, ok := task.Payload["prompt"].(string)
	if !ok || prompt == "" {
		// Try to construct a prompt based on the task type if one isn't explicitly provided.
		switch task.Type {
		case "ai.generate":
			prompt = "Please generate the requested code."
		case "ai.refactor":
			prompt = "Please refactor the selected code."
		case "ai.explain":
			prompt = "Please explain the selected code."
		default:
			return fmt.Errorf("missing 'prompt' in task payload")
		}
	}

	model, _ := task.Payload["model"].(string)

	// Extract editor state from payload
	editorState := make(map[string]any)
	if ctxMap, ok := task.Payload["context"].(map[string]any); ok {
		editorState = ctxMap
	} else {
		// Fallback for flat structure
		if file, ok := task.Payload["file"].(string); ok {
			editorState["file"] = file
		}
		if language, ok := task.Payload["language"].(string); ok {
			editorState["language"] = language
		}
		if selectedText, ok := task.Payload["selectedText"].(string); ok {
			editorState["selectedText"] = selectedText
		}
		if cursor, ok := task.Payload["cursor"].(float64); ok {
			editorState["cursor"] = cursor
		}
	}

	req := InferenceRequest{
		Prompt:      prompt,
		Model:       model,
		EditorState: editorState,
	}

	resp, err := r.AIRuntime.Infer(ctx, common.WorkspaceID(task.WorkspaceID), req)
	if err != nil {
		return fmt.Errorf("AI inference failed: %w", err)
	}

	task.Result = map[string]any{
		"text":     resp.Text,
		"model":    resp.Model,
		"provider": resp.ProviderName,
		"usage":    resp.Usage,
	}

	return nil
}

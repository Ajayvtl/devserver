package ai

import (
	"context"
	"fmt"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/Ajayvtl/devserver/internal/tasks"
)

// TaskRunner implements the tasks.Runner interface to handle AI intents (generate, refactor, explain)
// triggered via the Command Bus.
type TaskRunner struct {
	AIRuntime *Runtime
	Bus       events.Bus
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
	editorState := &EditorState{}
	if ctxMap, ok := task.Payload["context"].(map[string]any); ok {
		if file, ok := ctxMap["file"].(string); ok {
			editorState.ActiveFile = file
		}
		if language, ok := ctxMap["language"].(string); ok {
			editorState.Language = language
		}
		if selectedText, ok := ctxMap["selectedText"].(string); ok {
			editorState.SelectedText = selectedText
		}
		if cursor, ok := ctxMap["cursor"].(float64); ok {
			editorState.Cursor = cursor
		}
		if diags, ok := ctxMap["diagnostics"].([]any); ok {
			for _, item := range diags {
				if dMap, ok := item.(map[string]any); ok {
					var diag Diagnostic
					if f, ok := dMap["file"].(string); ok {
						diag.File = f
					}
					if l, ok := dMap["line"].(float64); ok {
						diag.Line = int(l)
					}
					if m, ok := dMap["message"].(string); ok {
						diag.Message = m
					}
					if s, ok := dMap["severity"].(string); ok {
						diag.Severity = s
					}
					editorState.Diagnostics = append(editorState.Diagnostics, diag)
				}
			}
		}
	} else {
		// Fallback for flat structure
		if file, ok := task.Payload["file"].(string); ok {
			editorState.ActiveFile = file
		}
		if language, ok := task.Payload["language"].(string); ok {
			editorState.Language = language
		}
		if selectedText, ok := task.Payload["selectedText"].(string); ok {
			editorState.SelectedText = selectedText
		}
		if cursor, ok := task.Payload["cursor"].(float64); ok {
			editorState.Cursor = cursor
		}
	}

	req := InferenceRequest{
		Prompt:      prompt,
		Model:       model,
		EditorState: editorState,
	}

	if sessionID, ok := task.Payload["sessionId"].(string); ok {
		req.SessionID = sessionID
	}

	stream, _ := task.Payload["stream"].(bool)
	if stream && r.Bus != nil {
		req.Stream = true
		req.StreamCallback = func(token string) {
			r.Bus.Publish(events.AIStreamToken, events.AIStreamTokenEvent{
				TaskID: task.ID,
				Token:  token,
			})
		}
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

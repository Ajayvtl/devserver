package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ajayvtl/devserver/internal/executor/legacy"
)

// AIProvider manages local AI models via Ollama using the legacy.Runtime.
type AIProvider struct {
	Runtime legacy.Runtime
	Meta    ProviderMetadata
}

func NewAIProvider(runtime legacy.Runtime) *AIProvider {
	return &AIProvider{
		Runtime: runtime,
		Meta: ProviderMetadata{
			Name:        "ollama",
			Description: "Local AI Runtime (Ollama)",
			Version:     "latest",
			Tags:        []string{"ai", "llm", "local"},
			Properties: map[string]any{
				"capabilities": map[string]bool{
					"models":     true,
					"embeddings": true,
					"streaming":  true,
				},
			},
		},
	}
}

func (p *AIProvider) Metadata() ProviderMetadata {
	return p.Meta
}

func (p *AIProvider) Detect(ctx context.Context) (bool, error) {
	return p.Runtime.DetectBinary(ctx, "ollama")
}

func (p *AIProvider) Version(ctx context.Context) (string, error) {
	installed, err := p.Detect(ctx)
	if !installed || err != nil {
		return "unknown", nil
	}
	out, err := p.Runtime.QueryVersion(ctx, "ollama", "--version")
	if err != nil {
		return "unknown", fmt.Errorf("failed to get ollama version: %w", err)
	}
	return out, nil
}

func (p *AIProvider) Health(ctx context.Context) error {
	installed, _ := p.Detect(ctx)
	if !installed {
		return fmt.Errorf("ollama is not installed")
	}
	// A basic health check could try to list models
	_, err := p.Runtime.Execute(ctx, "ollama", "list")
	if err != nil {
		return fmt.Errorf("ollama service is not responding: %w", err)
	}
	return nil
}

func (p *AIProvider) Status(ctx context.Context) (ProviderStatus, error) {
	installed, _ := p.Detect(ctx)
	if !installed {
		return StatusNotInstalled, nil
	}
	err := p.Health(ctx)
	if err != nil {
		return StatusStopped, nil
	}
	return StatusRunning, nil
}

func (p *AIProvider) Info(ctx context.Context) (ProviderInfo, error) {
	status, _ := p.Status(ctx)
	health := HealthUnknown
	if status == StatusStopped {
		health = HealthFailed
	} else if status == StatusNotInstalled {
		health = HealthNotApplicable
	} else if err := p.Health(ctx); err == nil {
		health = HealthHealthy
	} else {
		health = HealthFailed
	}

	version, _ := p.Version(ctx)
	installed, _ := p.Detect(ctx)

	if !installed {
		version = ""
	}

	return ProviderInfo{
		Name:      p.Meta.Name,
		Version:   version,
		Installed: installed,
		State: ProviderState{
			Status: status,
			Health: health,
		},
		Metrics:      ProviderMetrics{},
		Capabilities: p.Capabilities(),
	}, nil
}

func (p *AIProvider) Capabilities() ProviderCapabilities {
	return ProviderCapabilities{
		Detect:  true,
		Version: true,
		Health:  true,
		Status:  true,
		Start:   true,
		Stop:    true,
		Restart: true,
	}
}

func (p *AIProvider) Install(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *AIProvider) Update(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *AIProvider) Uninstall(ctx context.Context) error {
	return legacy.ErrUnsupportedOperation
}

func (p *AIProvider) Start(ctx context.Context) error {
	// Starting the ollama server in the background
	_, err := p.Runtime.Execute(ctx, "ollama", "serve")
	if err != nil {
		return fmt.Errorf("failed to start ollama: %w", err)
	}
	return nil
}

func (p *AIProvider) Stop(ctx context.Context) error {
	_, err := p.Runtime.Execute(ctx, "pkill", "-f", "ollama serve")
	if err != nil {
		return fmt.Errorf("failed to stop ollama: %w", err)
	}
	return nil
}

func (p *AIProvider) Restart(ctx context.Context) error {
	_ = p.Stop(ctx)
	return p.Start(ctx)
}

func (p *AIProvider) Configure(ctx context.Context, config map[string]any) error {
	return legacy.ErrUnsupportedOperation
}

func (p *AIProvider) Validate(ctx context.Context) error {
	return p.Health(ctx)
}

func (p *AIProvider) Generate(ctx context.Context, model string, prompt string) (string, error) {
	if model == "" {
		model = "llama2" // Default fallback model
	}
	out, err := p.Runtime.Execute(ctx, "ollama", "run", model, prompt)
	if err != nil {
		if out != nil {
			return "", fmt.Errorf("ollama inference failed: %w (stderr: %s)", err, out.Stderr)
		}
		return "", fmt.Errorf("ollama inference failed: %w", err)
	}
	return out.Stdout, nil
}

func (p *AIProvider) GenerateStream(ctx context.Context, model string, prompt string, onToken func(string)) error {
	if model == "" {
		model = "llama2"
	}

	reqBody, err := json.Marshal(map[string]any{
		"model":  model,
		"prompt": prompt,
		"stream": true,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/generate", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ollama API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		line := scanner.Text()
		var chunk struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}
		if err := json.Unmarshal([]byte(line), &chunk); err == nil {
			if chunk.Response != "" {
				onToken(chunk.Response)
			}
			if chunk.Done {
				break
			}
		}
	}

	return scanner.Err()
}

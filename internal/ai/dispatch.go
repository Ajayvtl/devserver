package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ajayvtl/devserver/internal/providers"
	"github.com/rs/zerolog"
)

var (
	ErrNoSuitableProvider = errors.New("no suitable AI provider found")
	ErrProviderNotRunning = errors.New("selected AI provider is not running")
	ErrModelNotFound      = errors.New("requested model not available")
)

// InferenceRequest represents a prompt to the LLM along with associated options.
type InferenceRequest struct {
	Prompt         string
	Model          string
	Temperature    float64
	MaxTokens      int
	Stream         bool
	SessionID      string
	StreamCallback func(string)
	EditorState    *EditorState
	Context        *AssembledContext
}

// InferenceResponse represents a discrete response from the LLM.
type InferenceResponse struct {
	Text         string
	Model        string
	ProviderName string
	Usage        TokenUsage
}

// TokenUsage tracks token consumption for the request.
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Dispatcher handles routing inference requests to the appropriate AI Provider.
type Dispatcher struct {
	log      zerolog.Logger
	managers *providers.Manager
}

// NewDispatcher creates a new AI inference Dispatcher.
func NewDispatcher(logger zerolog.Logger, pm *providers.Manager) *Dispatcher {
	return &Dispatcher{
		log:      logger.With().Str("component", "AIDispatcher").Logger(),
		managers: pm,
	}
}

// Dispatch routes the request to an available AI provider based on capabilities.
// It relies on ProviderMetadata.Properties to identify AI capability.
func (d *Dispatcher) Dispatch(ctx context.Context, req InferenceRequest) (*InferenceResponse, error) {
	provider, err := d.selectProvider(ctx, req)
	if err != nil {
		return nil, err
	}

	d.log.Info().Str("provider", provider.Metadata().Name).Str("model", req.Model).Msg("Dispatching inference request")

	// Check cancellation before invoking the provider.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	infProvider, ok := provider.(providers.InferenceProvider)
	if !ok {
		return nil, fmt.Errorf("provider %s does not support inference", provider.Metadata().Name)
	}

	fullPrompt := req.Prompt

	if req.Stream && req.StreamCallback != nil {
		err := infProvider.GenerateStream(ctx, req.Model, fullPrompt, req.StreamCallback)
		if err != nil {
			return nil, err
		}
		return &InferenceResponse{
			Text:         "", // Streamed, no full text
			Model:        req.Model,
			ProviderName: provider.Metadata().Name,
			Usage:        TokenUsage{},
		}, nil
	}

	out, err := infProvider.Generate(ctx, req.Model, fullPrompt)
	if err != nil {
		return nil, err
	}

	return &InferenceResponse{
		Text:         out,
		Model:        req.Model,
		ProviderName: provider.Metadata().Name,
		Usage: TokenUsage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		}, // Token counts not exposed by local CLI runtime abstraction
	}, nil
}

// selectProvider finds a provider that supports AI capabilities and matches the request requirements.
func (d *Dispatcher) selectProvider(ctx context.Context, req InferenceRequest) (providers.Provider, error) {
	if d.managers == nil {
		return nil, ErrNoSuitableProvider
	}

	for _, p := range d.managers.List() {
		meta := p.Metadata()
		if meta.Properties == nil {
			continue
		}

		caps, ok := meta.Properties["capabilities"].(map[string]bool)
		if !ok {
			continue
		}

		// Check if it's an AI provider
		isAI := false
		for _, tag := range meta.Tags {
			if tag == "ai" || tag == "llm" {
				isAI = true
				break
			}
		}

		if !isAI {
			continue
		}

		// Validate required capabilities
		if req.Stream && !caps["streaming"] {
			continue
		}

		status, _ := p.Status(ctx)
		if status != providers.StatusRunning {
			d.log.Debug().Str("provider", meta.Name).Msg("Provider found but not running")
			continue
		}

		return p, nil
	}

	return nil, ErrNoSuitableProvider
}

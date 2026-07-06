package editor

import (
	"fmt"
	"sync"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/rs/zerolog"
)

type proxyManagerImpl struct {
	log    zerolog.Logger
	routes map[common.WorkspaceID]string
	mu     sync.RWMutex
}

// NewProxyManager creates a new ProxyManager.
func NewProxyManager(logger zerolog.Logger) ProxyManager {
	return &proxyManagerImpl{
		log:    logger.With().Str("component", "ProxyManager").Logger(),
		routes: make(map[common.WorkspaceID]string),
	}
}

func (p *proxyManagerImpl) RegisterRoute(workspaceID common.WorkspaceID, targetURL string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.routes[workspaceID] = targetURL
	p.log.Debug().Str("workspace", string(workspaceID)).Str("target", targetURL).Msg("Proxy route registered")
	return nil
}

func (p *proxyManagerImpl) RemoveRoute(workspaceID common.WorkspaceID) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.routes[workspaceID]; !ok {
		return fmt.Errorf("proxy route not found for workspace %s", workspaceID)
	}
	delete(p.routes, workspaceID)
	p.log.Debug().Str("workspace", string(workspaceID)).Msg("Proxy route removed")
	return nil
}

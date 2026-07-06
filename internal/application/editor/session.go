package editor

import (
	"context"
	"fmt"
	"sync"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/rs/zerolog"
)

type sessionManagerImpl struct {
	log      zerolog.Logger
	sessions map[common.WorkspaceID]*WorkspaceSession
	mu       sync.RWMutex
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager(logger zerolog.Logger) SessionManager {
	return &sessionManagerImpl{
		log:      logger.With().Str("component", "SessionManager").Logger(),
		sessions: make(map[common.WorkspaceID]*WorkspaceSession),
	}
}

func (s *sessionManagerImpl) GetSession(ctx context.Context, workspaceID common.WorkspaceID) (*WorkspaceSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[workspaceID]
	if !ok {
		return nil, fmt.Errorf("session not found for workspace %s", workspaceID)
	}
	return session, nil
}

func (s *sessionManagerImpl) SaveSession(ctx context.Context, session *WorkspaceSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.WorkspaceID] = session
	s.log.Debug().Str("workspace", string(session.WorkspaceID)).Msg("Session saved")
	return nil
}

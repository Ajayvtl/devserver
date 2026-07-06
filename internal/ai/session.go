package ai

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// MessageRole defines who authored the message
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
)

// Message represents a single turn in a conversation
type Message struct {
	ID        string
	Role      MessageRole
	Content   string
	Timestamp time.Time
}

// Conversation maintains the history of a session
type Conversation struct {
	ID        string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SessionManager manages active conversations in memory
type SessionManager struct {
	mu            sync.RWMutex
	conversations map[string]*Conversation
}

// NewSessionManager creates a new SessionManager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		conversations: make(map[string]*Conversation),
	}
}

// GetOrCreate returns an existing conversation or creates a new one
func (sm *SessionManager) GetOrCreate(id string) *Conversation {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if id == "" {
		id = uuid.NewString()
	}

	conv, exists := sm.conversations[id]
	if !exists {
		conv = &Conversation{
			ID:        id,
			Messages:  make([]Message, 0),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		sm.conversations[id] = conv
	}

	return conv
}

// AddMessage appends a message to the conversation history
func (sm *SessionManager) AddMessage(sessionID string, role MessageRole, content string) *Message {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	conv, exists := sm.conversations[sessionID]
	if !exists {
		// Auto-create if not found to handle disconnected flows
		conv = &Conversation{
			ID:        sessionID,
			Messages:  make([]Message, 0),
			CreatedAt: time.Now().UTC(),
		}
		sm.conversations[sessionID] = conv
	}

	msg := Message{
		ID:        uuid.NewString(),
		Role:      role,
		Content:   content,
		Timestamp: time.Now().UTC(),
	}

	conv.Messages = append(conv.Messages, msg)
	conv.UpdatedAt = msg.Timestamp

	return &msg
}

// GetHistory returns all messages for a session
func (sm *SessionManager) GetHistory(sessionID string) []Message {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if conv, exists := sm.conversations[sessionID]; exists {
		// Return a copy to prevent data races
		msgs := make([]Message, len(conv.Messages))
		copy(msgs, conv.Messages)
		return msgs
	}

	return nil
}

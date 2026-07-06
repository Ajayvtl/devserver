package runtime

import (
	"time"
)

// Event represents a standardized envelope for all system events.
// This strictly replaces unstructured map[string]any payloads.
type Event struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Source  string    `json:"source"`
	Time    time.Time `json:"time"`
	Payload any       `json:"payload"`
}

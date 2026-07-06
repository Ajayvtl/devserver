package events

import "time"

func (b *eventBus) Publish(eventType EventType, payload any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	evt := Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
	}

	for _, ch := range b.subscribers[eventType] {
		select {
		case ch <- evt:
		default:
			// Non-blocking publish
		}
	}
}

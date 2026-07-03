package events

import "sync"

type Bus interface {
	Publish(eventType EventType, payload any)
	Subscribe(eventType EventType) Subscriber
}

type eventBus struct {
	mu          sync.RWMutex
	subscribers map[EventType][]chan Event
}

func NewBus() Bus {
	return &eventBus{
		subscribers: make(map[EventType][]chan Event),
	}
}

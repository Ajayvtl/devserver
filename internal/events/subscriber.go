package events

func (b *eventBus) Subscribe(eventType EventType) Subscriber {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan Event, 100)
	b.subscribers[eventType] = append(b.subscribers[eventType], ch)

	return ch
}

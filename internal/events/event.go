package events

import "time"

type Event struct {
	Type      EventType
	Payload   any
	Timestamp time.Time
}

type Subscriber <-chan Event

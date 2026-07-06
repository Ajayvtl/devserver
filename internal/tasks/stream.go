package tasks

import (
	"context"
	"net/http"
	"sync"

	"github.com/Ajayvtl/devserver/internal/events"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/rs/zerolog"
)

// EventStream bridges the EventBus to WebSocket clients.
// It replaces the task-specific TaskHub with a generalized channel
// that carries task, workspace, plugin, and AI events.
type EventStream struct {
	log     zerolog.Logger
	bus     events.Bus
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
	ctx     context.Context
	cancel  context.CancelFunc
	status  rt.Status
}

// NewEventStream creates a new stream that subscribes to the EventBus.
func NewEventStream(log zerolog.Logger, bus events.Bus) *EventStream {
	return &EventStream{
		log:     log,
		bus:     bus,
		clients: make(map[*websocket.Conn]struct{}),
		status:  rt.StatusStopped,
	}
}

func (s *EventStream) Name() string { return "tasks.EventStream" }

func (s *EventStream) Initialize(ctx context.Context) error {
	s.status = rt.StatusStarting
	return nil
}

func (s *EventStream) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.status = rt.StatusRunning
	go s.runInternal(s.ctx)
	return nil
}

func (s *EventStream) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	s.status = rt.StatusStopped
	return nil
}

func (s *EventStream) Status() rt.Status { return s.status }

func (s *EventStream) Health() rt.Health { return rt.HealthHealthy }

// wireEvent is the JSON shape sent to the frontend.
type wireEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// runInternal subscribes to all relevant event types and broadcasts them
// to connected WebSocket clients. Blocks until ctx is cancelled.
func (s *EventStream) runInternal(ctx context.Context) {
	// Subscribe to all event types we want to forward.
	topics := []events.EventType{
		events.TaskQueued,
		events.TaskStarted,
		events.TaskProgress,
		events.TaskCompleted,
		events.TaskFailed,
		events.TaskCancelled,
		events.TaskRolledBack,
		events.CommandSubmitted,
		events.CommandStarted,
		events.CommandCompleted,
		events.CommandFailed,
		events.CommandCancelled,
		events.CommandRolledBack,
		events.WorkspaceChanged,
		events.WorkspaceActivated,
		events.WorkspaceDeactivated,
		events.WorkspaceIndexed,
		events.ContextGenerated,
		events.AIStreamToken,
	}

	channels := make([]events.Subscriber, len(topics))
	for i, topic := range topics {
		channels[i] = s.bus.Subscribe(topic)
	}

	s.log.Info().Int("topics", len(topics)).Msg("event stream started")

	// Fan-in: merge all subscription channels.
	merged := make(chan wireEvent, 256)
	var wg sync.WaitGroup
	for i, ch := range channels {
		wg.Add(1)
		go func(topic events.EventType, sub events.Subscriber) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case evt, ok := <-sub:
					if !ok {
						return
					}
					merged <- wireEvent{
						Type:    string(topic),
						Payload: evt.Payload,
					}
				}
			}
		}(topics[i], ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	for {
		select {
		case <-ctx.Done():
			s.closeAll()
			return
		case evt, ok := <-merged:
			if !ok {
				return
			}
			s.broadcast(evt)
		}
	}
}

// HandleWS is the HTTP handler for WebSocket connections.
// Mount this at /ws/events to replace /ws/tasks.
func (s *EventStream) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		return
	}
	s.mu.Lock()
	s.clients[conn] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}()

	// Keep connection alive by reading (and discarding) client messages.
	ctx := r.Context()
	for {
		var ignore map[string]any
		if err := wsjson.Read(ctx, conn, &ignore); err != nil {
			return
		}
	}
}

func (s *EventStream) broadcast(evt wireEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for conn := range s.clients {
		_ = wsjson.Write(context.Background(), conn, evt)
	}
}

func (s *EventStream) closeAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for conn := range s.clients {
		_ = conn.Close(websocket.StatusNormalClosure, "")
		delete(s.clients, conn)
	}
}

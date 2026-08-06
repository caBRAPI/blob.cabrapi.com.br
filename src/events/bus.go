package events

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// EventType identifies the kind of internal event.
type EventType string

const (
	BlobUploaded       EventType = "blob.uploaded"
	BlobDeleted        EventType = "blob.deleted"
	BlobExpired        EventType = "blob.expired"
	MultipartCompleted EventType = "multipart.completed"
	SignedURLCreated   EventType = "signed_url.created"
)

// Event is an internal, in-process event carrying the context of an action.
type Event struct {
	Type      EventType
	BlobID    uuid.UUID
	Bucket    string
	Filename  string
	Size      int64
	Timestamp time.Time
	Data      map[string]interface{}
}

// Handler receives published events.
type Handler func(Event)

// Bus is a lightweight in-process pub/sub for internal events. It enables
// logging, metrics and future plugins without coupling callers.
type Bus struct {
	mu       sync.RWMutex
	handlers map[EventType][]Handler
}

// NewBus creates an empty event bus.
func NewBus() *Bus {
	return &Bus{handlers: make(map[EventType][]Handler)}
}

// Subscribe registers a handler for an event type and returns an unsubscribe func.
func (b *Bus) Subscribe(t EventType, h Handler) func() {
	b.mu.Lock()
	b.handlers[t] = append(b.handlers[t], h)
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		for i, existing := range b.handlers[t] {
			if &existing == &h {
				b.handlers[t] = append(b.handlers[t][:i], b.handlers[t][i+1:]...)
				break
			}
		}
	}
}

// Publish delivers the event to all handlers asynchronously.
func (b *Bus) Publish(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b.mu.RLock()
	handlers := make([]Handler, 0, len(b.handlers[e.Type]))
	handlers = append(handlers, b.handlers[e.Type]...)
	b.mu.RUnlock()

	for _, h := range handlers {
		go func(handler Handler, ev Event) {
			defer func() {
				if rec := recover(); rec != nil {
					// a panicking handler must not take the process down
					_ = rec
				}
			}()
			handler(ev)
		}(h, e)
	}
}

// Default is the process-wide event bus used by services.
var Default = NewBus()

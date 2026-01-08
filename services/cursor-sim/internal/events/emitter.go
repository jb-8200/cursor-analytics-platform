package events

import "sync"

// Handler processes events.
// Handlers are called synchronously when an event is emitted.
type Handler func(Event)

// Emitter broadcasts events to subscribers.
// Implementations must be thread-safe.
type Emitter interface {
	// Emit broadcasts an event to all subscribed handlers
	Emit(Event)

	// Subscribe registers a handler to receive events
	Subscribe(Handler)

	// Unsubscribe removes a handler (optional, not all implementations support this)
	Unsubscribe(Handler)
}

// MemoryEmitter is an in-memory event emitter with thread-safe subscription handling.
// Handlers are called synchronously in the order they were subscribed.
type MemoryEmitter struct {
	handlers []Handler
	mu       sync.RWMutex
}

// NewMemoryEmitter creates a new in-memory event emitter
func NewMemoryEmitter() *MemoryEmitter {
	return &MemoryEmitter{
		handlers: make([]Handler, 0),
	}
}

// Emit broadcasts an event to all subscribed handlers.
// Handlers are called synchronously in subscription order.
// Thread-safe: can be called concurrently from multiple goroutines.
func (e *MemoryEmitter) Emit(event Event) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, handler := range e.handlers {
		handler(event)
	}
}

// Subscribe registers a handler to receive all future events.
// Thread-safe: can be called concurrently from multiple goroutines.
// The same handler can be subscribed multiple times (will be called multiple times per event).
func (e *MemoryEmitter) Subscribe(handler Handler) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.handlers = append(e.handlers, handler)
}

// Unsubscribe is not implemented for MemoryEmitter.
// This is acceptable because handlers are typically subscribed at startup and never removed.
func (e *MemoryEmitter) Unsubscribe(handler Handler) {
	// Not implemented - handlers live for the lifetime of the emitter
	// This is acceptable for cursor-sim's use case where handlers are set up once
}

// NullEmitter discards all events and ignores all subscriptions.
// Use this in tests when you want to disable event emission.
type NullEmitter struct{}

// Emit discards the event (no-op)
func (e *NullEmitter) Emit(Event) {
	// Discard event
}

// Subscribe discards the handler (no-op)
func (e *NullEmitter) Subscribe(Handler) {
	// Discard handler
}

// Unsubscribe is a no-op
func (e *NullEmitter) Unsubscribe(Handler) {
	// No-op
}

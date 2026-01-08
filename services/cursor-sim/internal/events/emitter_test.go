package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryEmitter_Emit_SingleSubscriber(t *testing.T) {
	emitter := NewMemoryEmitter()

	var received Event
	emitter.Subscribe(func(e Event) {
		received = e
	})

	event := ProgressEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeProgress,
			Time:      time.Now(),
		},
		Phase:   "test",
		Current: 5,
		Total:   10,
	}

	emitter.Emit(event)

	require.NotNil(t, received)
	assert.Equal(t, EventTypeProgress, received.Type())

	// Cast to ProgressEvent to check fields
	progressEvent, ok := received.(ProgressEvent)
	require.True(t, ok)
	assert.Equal(t, "test", progressEvent.Phase)
	assert.Equal(t, 5, progressEvent.Current)
	assert.Equal(t, 10, progressEvent.Total)
}

func TestMemoryEmitter_Emit_MultipleSubscribers(t *testing.T) {
	emitter := NewMemoryEmitter()

	count := 0
	handler := func(e Event) { count++ }

	emitter.Subscribe(handler)
	emitter.Subscribe(handler)
	emitter.Subscribe(handler)

	emitter.Emit(PhaseStartEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypePhaseStart,
			Time:      time.Now(),
		},
		Phase:   "test",
		Message: "Starting test",
	})

	assert.Equal(t, 3, count)
}

func TestMemoryEmitter_Emit_NoSubscribers(t *testing.T) {
	emitter := NewMemoryEmitter()

	// Should not panic when emitting with no subscribers
	assert.NotPanics(t, func() {
		emitter.Emit(ProgressEvent{
			BaseEvent: BaseEvent{
				EventType: EventTypeProgress,
				Time:      time.Now(),
			},
			Phase:   "test",
			Current: 1,
			Total:   10,
		})
	})
}

func TestMemoryEmitter_Subscribe_MultipleHandlers(t *testing.T) {
	emitter := NewMemoryEmitter()

	handler1Called := false
	handler2Called := false

	emitter.Subscribe(func(e Event) { handler1Called = true })
	emitter.Subscribe(func(e Event) { handler2Called = true })

	emitter.Emit(PhaseCompleteEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypePhaseComplete,
			Time:      time.Now(),
		},
		Phase:   "test",
		Message: "Test complete",
		Success: true,
	})

	assert.True(t, handler1Called)
	assert.True(t, handler2Called)
}

func TestMemoryEmitter_Emit_DifferentEventTypes(t *testing.T) {
	emitter := NewMemoryEmitter()

	var receivedEvents []Event
	emitter.Subscribe(func(e Event) {
		receivedEvents = append(receivedEvents, e)
	})

	// Emit different event types
	emitter.Emit(PhaseStartEvent{
		BaseEvent: BaseEvent{EventType: EventTypePhaseStart, Time: time.Now()},
		Phase:     "phase1",
	})
	emitter.Emit(ProgressEvent{
		BaseEvent: BaseEvent{EventType: EventTypeProgress, Time: time.Now()},
		Phase:     "phase1",
		Current:   1,
		Total:     10,
	})
	emitter.Emit(PhaseCompleteEvent{
		BaseEvent: BaseEvent{EventType: EventTypePhaseComplete, Time: time.Now()},
		Phase:     "phase1",
		Success:   true,
	})
	emitter.Emit(WarningEvent{
		BaseEvent: BaseEvent{EventType: EventTypeWarning, Time: time.Now()},
		Message:   "Test warning",
	})

	require.Equal(t, 4, len(receivedEvents))
	assert.Equal(t, EventTypePhaseStart, receivedEvents[0].Type())
	assert.Equal(t, EventTypeProgress, receivedEvents[1].Type())
	assert.Equal(t, EventTypePhaseComplete, receivedEvents[2].Type())
	assert.Equal(t, EventTypeWarning, receivedEvents[3].Type())
}

func TestMemoryEmitter_ThreadSafety(t *testing.T) {
	emitter := NewMemoryEmitter()

	var wg sync.WaitGroup
	eventCount := 0
	var mu sync.Mutex

	// Add subscriber that counts events
	emitter.Subscribe(func(e Event) {
		mu.Lock()
		eventCount++
		mu.Unlock()
	})

	// Emit events from multiple goroutines concurrently
	numGoroutines := 10
	eventsPerGoroutine := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				emitter.Emit(ProgressEvent{
					BaseEvent: BaseEvent{
						EventType: EventTypeProgress,
						Time:      time.Now(),
					},
					Phase:   "concurrent_test",
					Current: j,
					Total:   eventsPerGoroutine,
				})
			}
		}()
	}

	wg.Wait()

	// Should have received all events
	mu.Lock()
	assert.Equal(t, numGoroutines*eventsPerGoroutine, eventCount)
	mu.Unlock()
}

func TestMemoryEmitter_SubscribeBeforeEmit(t *testing.T) {
	emitter := NewMemoryEmitter()

	firstHandlerCalled := false
	secondHandlerCalled := false

	// Subscribe two handlers before emitting
	emitter.Subscribe(func(e Event) {
		firstHandlerCalled = true
	})
	emitter.Subscribe(func(e Event) {
		secondHandlerCalled = true
	})

	// Emit - both handlers should be called
	emitter.Emit(ProgressEvent{
		BaseEvent: BaseEvent{EventType: EventTypeProgress, Time: time.Now()},
		Phase:     "test",
		Current:   1,
		Total:     2,
	})

	assert.True(t, firstHandlerCalled)
	assert.True(t, secondHandlerCalled)
}

func TestNullEmitter_Emit(t *testing.T) {
	emitter := &NullEmitter{}

	// Should not panic
	assert.NotPanics(t, func() {
		emitter.Emit(ProgressEvent{
			BaseEvent: BaseEvent{EventType: EventTypeProgress, Time: time.Now()},
			Phase:     "test",
			Current:   5,
			Total:     10,
		})
	})
}

func TestNullEmitter_Subscribe(t *testing.T) {
	emitter := &NullEmitter{}

	handlerCalled := false

	// Should not panic
	assert.NotPanics(t, func() {
		emitter.Subscribe(func(e Event) {
			handlerCalled = true
		})
	})

	// Emit should not call the handler
	emitter.Emit(ProgressEvent{
		BaseEvent: BaseEvent{EventType: EventTypeProgress, Time: time.Now()},
		Phase:     "test",
		Current:   1,
		Total:     10,
	})

	assert.False(t, handlerCalled, "NullEmitter should not call handlers")
}

func TestNullEmitter_Unsubscribe(t *testing.T) {
	emitter := &NullEmitter{}

	// Should not panic
	assert.NotPanics(t, func() {
		emitter.Unsubscribe(func(e Event) {})
	})
}

func TestNewMemoryEmitter_InitialState(t *testing.T) {
	emitter := NewMemoryEmitter()

	assert.NotNil(t, emitter)
	assert.NotNil(t, emitter.handlers)
	assert.Equal(t, 0, len(emitter.handlers))
}

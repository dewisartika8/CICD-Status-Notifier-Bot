package events

// EventDispatcher handles domain events
type EventDispatcher interface {
	Dispatch(event DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
}

// EventHandler handles specific domain events
type EventHandler interface {
	Handle(event DomainEvent) error
	CanHandle(eventType string) bool
}

// SimpleEventDispatcher provides basic event dispatching
type SimpleEventDispatcher struct {
	handlers map[string][]EventHandler
}

// NewSimpleEventDispatcher creates a new simple event dispatcher
func NewSimpleEventDispatcher() *SimpleEventDispatcher {
	return &SimpleEventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// Dispatch dispatches an event to all registered handlers
func (d *SimpleEventDispatcher) Dispatch(event DomainEvent) error {
	if handlers, exists := d.handlers[event.EventType()]; exists {
		for _, handler := range handlers {
			if err := handler.Handle(event); err != nil {
				return err
			}
		}
	}
	return nil
}

// Subscribe adds a handler for a specific event type
func (d *SimpleEventDispatcher) Subscribe(eventType string, handler EventHandler) error {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
	return nil
}

// Unsubscribe removes a handler for a specific event type
func (d *SimpleEventDispatcher) Unsubscribe(eventType string, handler EventHandler) error {
	if handlers, exists := d.handlers[eventType]; exists {
		for i, h := range handlers {
			if h == handler {
				d.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
	return nil
}

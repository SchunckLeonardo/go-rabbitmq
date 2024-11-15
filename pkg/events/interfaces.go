package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() any
}

type EventHandlerInterface interface {
	// Handle Execute function to run on event
	Handle(event EventInterface, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	// Register Create a new event
	Register(eventName string, handler EventHandlerInterface) error
	// Dispatch Send event to message broker
	Dispatch(event EventInterface) error
	// Remove Delete an event
	Remove(eventName string, handler EventHandlerInterface) error
	// Has If event exists
	Has(eventName string, handler EventHandlerInterface) bool
	// Clear Clean all events
	Clear()
}

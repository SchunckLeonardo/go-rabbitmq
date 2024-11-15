package events

import (
	"errors"
	"sync"
)

var (
	ErrHandlerAlreadyRegistered = errors.New("handler already registered")
	ErrHandlerNotExists         = errors.New("handler not exists")
)

type Handles map[string][]EventHandlerInterface

type EventDispatcher struct {
	handlers Handles
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{handlers: make(Handles)}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	var wg sync.WaitGroup
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		wg.Add(len(ed.handlers[event.GetName()]))
		for _, h := range handlers {
			go h.Handle(event, &wg)
		}
		wg.Wait()
	}
	return nil
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	isHandlerExists := ed.Has(eventName, handler)
	if !isHandlerExists {
		return ErrHandlerNotExists
	}

	for i, h := range ed.handlers[eventName] {
		if h == handler {
			ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
			break
		}
	}
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *EventDispatcher) Clear() {
	ed.handlers = Handles{}
}

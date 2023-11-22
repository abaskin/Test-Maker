package eventbus

import (
	"log"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/pkg/errors"
)

type EventHandlerFunc func(string, interface{})

type EventBus struct {
	handlers map[string][]*EventHandler
	queue    *goconcurrentqueue.FIFO
}

type EventHandler struct {
	eventFunc EventHandlerFunc
	async     bool
}

type Event struct {
	eventID   string
	eventData interface{}
}

func New() *EventBus {
	e := &EventBus{
		handlers: make(map[string][]*EventHandler),
		queue:    goconcurrentqueue.NewFIFO(),
	}
	go func() {
		evnt, err := e.queue.DequeueOrWaitForNextElement()
		if err != nil {
			log.Println(errors.Wrap(err, "eventBus"))
		}
		event := evnt.(*Event)
		handlers, found := e.handlers[event.eventID]
		if found {
			for _, handler := range handlers {
				if handler.async {
					go handler.eventFunc(event.eventID, event.eventData)
					continue
				}
				handler.eventFunc(event.eventID, event.eventData)
			}
		}
	}()
	return e
}

func (e *EventBus) AddHandler(eventID string, eventFunc EventHandlerFunc,
	async bool) error {
	if _, found := e.handlers[eventID]; !found {
		e.handlers[eventID] = make([]*EventHandler, 0)
	}
	e.handlers[eventID] = append(e.handlers[eventID],
		&EventHandler{
			eventFunc: eventFunc,
			async:     async,
		})
	return nil
}

func (e *EventBus) PostEvent(eventID string, data interface{}) error {
	e.queue.Enqueue(
		&Event{
			eventID:   eventID,
			eventData: data,
		})
	return nil
}

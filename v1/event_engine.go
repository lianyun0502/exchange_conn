package exchange_conn

import "log"

type Event struct{
	Name string
	Data interface{}
	Handler func()
	IsBlock bool
}


type EventEngine struct {
	eventQueue chan *Event
	StopSignal chan struct{}
}

func (e *EventEngine) AddEvent(event *Event) {
	e.eventQueue <- event
}

func (e *EventEngine) Luanch() {
	e.eventQueue = make(chan *Event, 100)
	go func ()  {
		for {
			event := <- e.eventQueue
			log.Printf("Event: name=%s", event.Name)
			if event.Name == "exit" {
				return
			}
			if event.IsBlock {
				event.Handler()
			} else {
				go event.Handler()
			}
		}
	}()
}

func (e *EventEngine) Stop() {
	e.eventQueue <- &Event{Name: "exit"}
	e.StopSignal <- struct{}{}
}
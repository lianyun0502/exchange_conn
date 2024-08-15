package exchange_conn

import log "github.com/sirupsen/logrus"

type Event struct {
	Name    string
	Data    interface{}
	Handler func()
	IsBlock bool
}

type EventEngine struct {
	eventQueue chan *Event
	StopSignal chan struct{}
	Events     map[string][]func()
	Logger     *log.Logger
}

func NewEventEngine() *EventEngine {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05.000000", FullTimestamp: true})
	return &EventEngine{
		StopSignal: make(chan struct{}),
		Events:     make(map[string][]func()),
		Logger:     logger,
	}
}

func (e *EventEngine) AddEvent(event *Event) {
	e.eventQueue <- event
}

func (e *EventEngine) Luanch() {
	e.eventQueue = make(chan *Event, 100)
	go func() {
		for {
			event := <-e.eventQueue
			e.Logger.WithFields(log.Fields{"event": event.Name}).Info("Event")
			if event.Name == "exit" {
				e.StopSignal <- struct{}{}
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
}

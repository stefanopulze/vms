package event

import "time"

func NewManager() *Manager {
	return &Manager{
		eventQueue:  make(chan Event, 100),
		subscribers: make(map[Type][]Fn),
	}
}

type Fn = func(event Event)

type Event struct {
	Type      Type
	Timestamp time.Time
	Payload   interface{}
}

type Manager struct {
	eventQueue  chan Event
	subscribers map[Type][]Fn
}

func (m *Manager) Start() {
	go func() {
		for event := range m.eventQueue {
			for _, handler := range m.subscribers[event.Type] {
				handler(event)
			}
		}
	}()
}

func (m *Manager) Close() {
	close(m.eventQueue)
}

package event

import "time"

type Publisher interface {
	Publish(event Type, payload interface{})
}

func (m *Manager) Publish(event Type, payload interface{}) {
	m.eventQueue <- Event{
		Type:      event,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

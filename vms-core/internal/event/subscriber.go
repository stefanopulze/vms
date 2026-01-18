package event

type Subscriber interface {
	Subscribe(event Type, handler Fn)
}

func (m *Manager) Subscribe(event Type, handler Fn) {
	m.subscribers[event] = append(m.subscribers[event], handler)
}

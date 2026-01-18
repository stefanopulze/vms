package exporter

import (
	"fmt"
	"log/slog"
	"vms-core/internal/voltronic"
)

type Client interface {
	Name() string
	Close() error
	GeneralStatus(data *voltronic.DeviceGeneralStatus, mode string) error
}

func NewMultiple() *Multiple {
	return &Multiple{
		clients: make([]Client, 0),
	}
}

type Multiple struct {
	clients []Client
}

func (m *Multiple) GeneralStatus(data *voltronic.DeviceGeneralStatus, mode string) error {
	var err error
	for _, c := range m.clients {
		if err = c.GeneralStatus(data, mode); err != nil {
			slog.Error(fmt.Sprintf("exporter %s: %v", c.Name(), err))
		}
	}

	return nil
}

func (m *Multiple) Name() string {
	return "multiple"
}

func (m *Multiple) AddExporter(c Client) {
	m.clients = append(m.clients, c)
}

func (m *Multiple) Close() error {
	for _, c := range m.clients {
		_ = c.Close()
	}
	return nil
}

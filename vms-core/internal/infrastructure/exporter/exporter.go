package exporter

import (
	"vms-core/internal/voltronic"
)

type Client interface {
	Name() string
	Close() error
	GeneralStatus(data *voltronic.DeviceGeneralStatus, mode string) error
}

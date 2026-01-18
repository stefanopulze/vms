package service

import (
	"fmt"
	"log/slog"
	"vms-core/internal/cache"
	"vms-core/internal/event"
	"vms-core/internal/infrastructure/exporter"
	"vms-core/internal/voltronic"
)

func NewExporter(i *voltronic.Client, e exporter.Client, em event.Publisher, qs *cache.QuerySnapshot) *Exporter {
	return &Exporter{
		inverter:      i,
		exporter:      e,
		events:        em,
		querySnapshot: qs,
	}
}

type Exporter struct {
	exporter      exporter.Client
	inverter      *voltronic.Client
	events        event.Publisher
	querySnapshot *cache.QuerySnapshot
}

func (e Exporter) ReadStatusInformation() {
	pigs, err := e.inverter.QueryPIGS()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to query PIGS: %s", err.Error()))
		return
	}
	e.querySnapshot.SetGeneralStatus(pigs)

	mode := "n.d"
	mr, err := e.inverter.QueryMode()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to query mode: %s", err.Error()))
	} else {
		mode = mr.Mode
	}

	warnings, err := e.inverter.QueryWarning()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to query warning: %s", err.Error()))
	} else {

	}
	e.querySnapshot.SetWarnings(warnings)

	if err = e.exporter.GeneralStatus(pigs, mode); err != nil {
		slog.Error("cannot export inverter status", slog.Any("error", err))
	}

}

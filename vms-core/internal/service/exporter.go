package service

import (
	"fmt"
	"log/slog"
	"vms-core/internal/cache"
	"vms-core/internal/infrastructure/exporter"
	"vms-core/internal/voltronic"
)

func NewScheduledCommands(i *voltronic.Client, e exporter.Client, qs *cache.QuerySnapshot, wm *WarningMonitor) *ScheduledCommands {
	return &ScheduledCommands{
		inverter:       i,
		exporter:       e,
		querySnapshot:  qs,
		warningMonitor: wm,
	}
}

type ScheduledCommands struct {
	exporter       exporter.Client
	inverter       *voltronic.Client
	querySnapshot  *cache.QuerySnapshot
	warningMonitor *WarningMonitor
}

func (e ScheduledCommands) Read() {
	pigs, err := e.inverter.QueryPIGS()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to query PIGS: %s", err.Error()))
	} else {
		e.querySnapshot.SetGeneralStatus(pigs)
	}

	mode := "n.d"
	mr, err := e.inverter.QueryMode()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to query mode: %s", err.Error()))
	} else {
		mode = mr.Mode
	}
	e.querySnapshot.SetMode(mode)

	warnings, err := e.inverter.QueryWarning()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to query warning: %s", err.Error()))
	} else {
		e.querySnapshot.SetWarnings(warnings)
	}

	if err = e.exporter.GeneralStatus(pigs, mode); err != nil {
		slog.Error("cannot export inverter status", slog.Any("error", err))
	}

	go e.warningMonitor.Check(pigs, mode, warnings)
}

package service

import (
	"context"
	"fmt"
	"vms-core/internal/notifier"
	"vms-core/internal/store"
	"vms-core/internal/voltronic"
)

type WarningMonitor struct {
	notifier         notifier.Notifier
	store            store.Store
	batteryThreshold []int
	batteryNotified  map[int]bool
}

func NewWarningMonitor(n notifier.Notifier, s store.Store) *WarningMonitor {
	return &WarningMonitor{
		notifier:         n,
		store:            s,
		batteryThreshold: []int{20, 30, 50, 80},
		batteryNotified:  make(map[int]bool),
	}
}

func (w *WarningMonitor) Check(pigs *voltronic.DeviceGeneralStatus, mode string, warnings *voltronic.DeviceWarning) {
	// Check for Mode Change
	var lastMode string
	if err := w.store.Load("mode", &lastMode); err == nil {
		if lastMode != mode {
			_ = w.notifier.Send(context.Background(), fmt.Sprintf("Mode changed from %s to %s", lastMode, mode))
			_ = w.store.Save("mode", mode)
		}
	} else {
		// First run or error loading, just save current mode
		_ = w.store.Save("mode", mode)
	}

	w.checkBatteryLevel(pigs.BatteryCapacity)
}

func (w *WarningMonitor) checkBatteryLevel(pct int) {
	for _, threshold := range w.batteryThreshold {
		if pct <= threshold && !w.batteryNotified[threshold] {
			msg := fmt.Sprintf("Battery is less than %d%%\nActual: %d%%", threshold, pct)
			w.batteryNotified[threshold] = true
			_ = w.notifier.Send(context.Background(), msg)
			return
		}

		if pct >= threshold+5 && w.batteryNotified[threshold] {
			w.batteryNotified[threshold] = false
		}
	}
}

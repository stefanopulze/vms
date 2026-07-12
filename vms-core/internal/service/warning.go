package service

import (
	"context"
	"fmt"
	"strings"
	"vms-core/internal/humanize"
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

func (w *WarningMonitor) Check(piri *voltronic.DeviceRatingInfo, pigs *voltronic.DeviceGeneralStatus, mode string, _ *voltronic.DeviceWarning) {
	if piri != nil {
		w.checkOutputSourcePriority(piri, mode)
	}
	if pigs != nil {
		w.checkBatteryLevel(pigs.BatteryCapacity)
	}
}

func (w *WarningMonitor) checkOutputSourcePriority(piri *voltronic.DeviceRatingInfo, mode string) {
	var lastMode string
	if err := w.store.Load("mode", &lastMode); err != nil {
		_ = w.store.Save("mode", mode)
		return
	}

	if lastMode != mode {
		_ = w.notifier.Send(context.Background(), fmt.Sprintf(
			"Mode %s changed to %s",
			strings.ToUpper(piri.OutputSourcePriorityEnum()),
			humanize.Mode(mode),
		))
		_ = w.store.Save("mode", mode)
	}
}

func (w *WarningMonitor) checkBatteryLevel(pct int) {
	// batteryThreshold is ascending, so the first breached entry is the most severe.
	lowest := -1
	charging := false

	for _, threshold := range w.batteryThreshold {
		switch {
		// crossed below a threshold while discharging: mark it, remember the lowest
		case pct <= threshold && !w.batteryNotified[threshold]:
			w.batteryNotified[threshold] = true
			if lowest == -1 {
				lowest = threshold
			}

		// recovered above a threshold (with 5% hysteresis): reset and flag recovery
		case pct >= threshold+5 && w.batteryNotified[threshold]:
			w.batteryNotified[threshold] = false
			charging = true
		}
	}

	// a single alert for the most severe threshold crossed this tick
	if lowest != -1 {
		_ = w.notifier.Send(context.Background(), fmt.Sprintf("Battery is less than %d%%", lowest))
	}

	// a single recovery notification even if several thresholds cleared at once
	if charging {
		_ = w.notifier.Send(context.Background(), fmt.Sprintf("🎉 Battery is charging %d%%", pct))
	}
}

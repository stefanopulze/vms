package service

import (
	"context"
	"fmt"
	"strings"
	"vms-core/internal/notifier"
	"vms-core/internal/store"
	"vms-core/internal/utils"
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

func (w *WarningMonitor) Check(piri *voltronic.DeviceRatingInfo, pigs *voltronic.DeviceGeneralStatus, mode string, warnings *voltronic.DeviceWarning) {
	w.checkOutputSourcePriority(piri, mode)
	w.checkBatteryLevel(pigs.BatteryCapacity)
}

func (w *WarningMonitor) checkOutputSourcePriority(piri *voltronic.DeviceRatingInfo, mode string) {
	var lastMode string
	if err := w.store.Load("mode", &lastMode); err != nil {
		_ = w.store.Save("mode", mode)
		return
	}

	if lastMode != mode {
		_ = w.notifier.Send(context.Background(), fmt.Sprintf(
			"Mode %s changed from %s to %s",
			strings.ToUpper(piri.OutputSourcePriorityEnum()),
			utils.ModeToHuman(lastMode),
			utils.ModeToHuman(mode),
		))
		_ = w.store.Save("mode", mode)
	}
}

func (w *WarningMonitor) checkBatteryLevel(pct int) {
	for _, threshold := range w.batteryThreshold {
		if pct <= threshold && !w.batteryNotified[threshold] {
			msg := fmt.Sprintf("Battery is less than %d%%", threshold)
			w.batteryNotified[threshold] = true
			_ = w.notifier.Send(context.Background(), msg)
			return
		}

		if pct >= threshold+5 && w.batteryNotified[threshold] {
			w.batteryNotified[threshold] = false
		}
	}
}

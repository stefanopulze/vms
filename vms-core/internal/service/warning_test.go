package service

import (
	"context"
	"testing"
	"vms-core/internal/voltronic"
)

// Mock Notifier
type MockNotifier struct {
	Messages []string
}

func (m *MockNotifier) Name() string {
	return "mock"
}

func (m *MockNotifier) Send(_ context.Context, message string) error {
	m.Messages = append(m.Messages, message)
	return nil
}

// Mock Store
type MockStore struct {
	Data map[string]interface{}
}

func NewMockStore() *MockStore {
	return &MockStore{
		Data: make(map[string]interface{}),
	}
}

func (m *MockStore) Save(key string, value interface{}) error {
	m.Data[key] = value
	return nil
}

func (m *MockStore) Load(key string, value interface{}) error {
	val, ok := m.Data[key]
	if !ok {
		return context.DeadlineExceeded // any "not found" error
	}
	if vStr, ok := val.(string); ok {
		if destStr, ok := value.(*string); ok {
			*destStr = vStr
			return nil
		}
	}
	return nil
}

func newMonitor() (*WarningMonitor, *MockNotifier, *MockStore) {
	notifier := &MockNotifier{}
	store := NewMockStore()
	return NewWarningMonitor(notifier, store), notifier, store
}

// Check must not panic when the inverter queries failed and nil snapshots are passed.
func TestWarningMonitor_Check_NilInputs(t *testing.T) {
	wm, notifier, _ := newMonitor()

	wm.Check(nil, nil, "line_mode", nil)

	if len(notifier.Messages) != 0 {
		t.Errorf("expected no notifications for nil inputs, got %v", notifier.Messages)
	}
}

func TestWarningMonitor_BatteryLevel(t *testing.T) {
	wm, notifier, _ := newMonitor()

	check := func(pct int) {
		notifier.Messages = nil
		wm.checkBatteryLevel(pct)
	}

	// Healthy battery: no alert.
	check(90)
	if len(notifier.Messages) != 0 {
		t.Fatalf("expected no alert at 90%%, got %v", notifier.Messages)
	}

	// Drops just below 80 -> single alert for the crossed threshold.
	check(79)
	if len(notifier.Messages) != 1 || notifier.Messages[0] != "Battery is less than 80%" {
		t.Fatalf("expected single '< 80%%' alert, got %v", notifier.Messages)
	}

	// Still below 80 but not recovered -> no repeat alert.
	check(79)
	if len(notifier.Messages) != 0 {
		t.Fatalf("expected no repeat alert while staying below 80%%, got %v", notifier.Messages)
	}

	// Continues to drop past 50 -> alert only for the newly crossed threshold.
	check(49)
	if len(notifier.Messages) != 1 || notifier.Messages[0] != "Battery is less than 50%" {
		t.Fatalf("expected single '< 50%%' alert, got %v", notifier.Messages)
	}
}

// A large single drop should report only the most severe threshold, not one per band.
func TestWarningMonitor_BatteryLevel_SingleAlertPerTick(t *testing.T) {
	wm, notifier, _ := newMonitor()

	wm.checkBatteryLevel(10) // below 80, 50, 30 and 20 at once

	if len(notifier.Messages) != 1 {
		t.Fatalf("expected a single alert for the most severe threshold, got %v", notifier.Messages)
	}
	if notifier.Messages[0] != "Battery is less than 20%" {
		t.Errorf("expected '< 20%%' (most severe), got %q", notifier.Messages[0])
	}
}

// Recovery (with hysteresis) resets thresholds and emits one charging notification.
func TestWarningMonitor_BatteryLevel_Recovery(t *testing.T) {
	wm, notifier, _ := newMonitor()

	wm.checkBatteryLevel(10) // trip every threshold

	notifier.Messages = nil
	wm.checkBatteryLevel(90) // recovers above all thresholds + hysteresis
	if len(notifier.Messages) != 1 {
		t.Fatalf("expected a single charging notification, got %v", notifier.Messages)
	}
	if notifier.Messages[0] != "🎉 Battery is charging 90%" {
		t.Errorf("unexpected charging message: %q", notifier.Messages[0])
	}

	// After recovery, dropping again must alert once more.
	notifier.Messages = nil
	wm.checkBatteryLevel(79)
	if len(notifier.Messages) != 1 || notifier.Messages[0] != "Battery is less than 80%" {
		t.Errorf("expected '< 80%%' alert after recovery, got %v", notifier.Messages)
	}
}

// Hysteresis: rising just above a threshold (but within 5%%) does not count as recovery.
func TestWarningMonitor_BatteryLevel_Hysteresis(t *testing.T) {
	wm, notifier, _ := newMonitor()

	wm.checkBatteryLevel(19) // trips the 20%% threshold

	notifier.Messages = nil
	wm.checkBatteryLevel(23) // above 20 but below 20+5 -> not recovered
	if len(notifier.Messages) != 0 {
		t.Fatalf("expected no recovery within hysteresis band, got %v", notifier.Messages)
	}

	wm.checkBatteryLevel(25) // reaches 20+5 -> recovery
	if len(notifier.Messages) != 1 || notifier.Messages[0] != "🎉 Battery is charging 25%" {
		t.Errorf("expected charging notification at hysteresis boundary, got %v", notifier.Messages)
	}
}

func TestWarningMonitor_ModeChange(t *testing.T) {
	wm, notifier, store := newMonitor()

	piri := &voltronic.DeviceRatingInfo{OutputSourcePriority: 2} // "sbu"
	healthy := &voltronic.DeviceGeneralStatus{BatteryCapacity: 90}

	// 1. First observation: unknown previous mode -> store it, no notification.
	wm.Check(piri, healthy, "line_mode", nil)
	if len(notifier.Messages) != 0 {
		t.Errorf("expected no notification on initial mode set, got %v", notifier.Messages)
	}
	if store.Data["mode"] != "line_mode" {
		t.Errorf("expected mode saved as line_mode, got %v", store.Data["mode"])
	}

	// 2. Same mode -> no notification.
	wm.Check(piri, healthy, "line_mode", nil)
	if len(notifier.Messages) != 0 {
		t.Errorf("expected no notification on unchanged mode, got %v", notifier.Messages)
	}

	// 3. Mode change -> one notification, updated store.
	wm.Check(piri, healthy, "battery_mode", nil)
	if len(notifier.Messages) != 1 {
		t.Fatalf("expected 1 notification on mode change, got %d (%v)", len(notifier.Messages), notifier.Messages)
	}
	if want := "Mode SBU changed to battery"; notifier.Messages[0] != want {
		t.Errorf("expected notification %q, got %q", want, notifier.Messages[0])
	}
	if store.Data["mode"] != "battery_mode" {
		t.Errorf("expected mode updated to battery_mode, got %v", store.Data["mode"])
	}
}

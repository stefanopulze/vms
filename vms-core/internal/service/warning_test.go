package service

import (
	"context"
	"strings"
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

func (m *MockNotifier) Send(ctx context.Context, message string) error {
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
		return context.DeadlineExceeded // Just an error
	}
	// Similar simple hack for mock as in real store
	if vStr, ok := val.(string); ok {
		if destStr, ok := value.(*string); ok {
			*destStr = vStr
			return nil
		}
	}
	return nil
}

func TestWarningMonitor_Check(t *testing.T) {
	mockNotifier := &MockNotifier{}
	mockStore := NewMockStore()
	wm := NewWarningMonitor(mockNotifier, mockStore)

	tests := []struct {
		name     string
		pigs     *voltronic.DeviceGeneralStatus
		warnings *voltronic.DeviceWarning
		wantMsg  string
	}{
		{
			name: "No warnings",
			pigs: &voltronic.DeviceGeneralStatus{},
			warnings: &voltronic.DeviceWarning{
				OverTemperature: false,
				OverLoad:        false,
			},
			wantMsg: "",
		},
		{
			name: "Over Temperature",
			pigs: &voltronic.DeviceGeneralStatus{},
			warnings: &voltronic.DeviceWarning{
				OverTemperature: true,
			},
			wantMsg: "Over Temperature",
		},
		{
			name: "Over Load",
			pigs: &voltronic.DeviceGeneralStatus{},
			warnings: &voltronic.DeviceWarning{
				OverLoad: true,
			},
			wantMsg: "Over Load",
		},
		{
			name: "High Heat Sink Temperature",
			pigs: &voltronic.DeviceGeneralStatus{
				InverterHeatSinkTemperature: 85,
			},
			warnings: &voltronic.DeviceWarning{},
			wantMsg:  "High Heat Sink Temperature: 85",
		},
		{
			name: "Multiple Warnings",
			pigs: &voltronic.DeviceGeneralStatus{
				InverterHeatSinkTemperature: 90,
			},
			warnings: &voltronic.DeviceWarning{
				OverTemperature:   true,
				BatteryVoltageLow: true,
			},
			wantMsg: "Over Temperature", // Just checking containment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNotifier.Messages = []string{} // Reset messages
			// Mode is currently not used in logic, pass empty
			wm.Check(tt.pigs, "", tt.warnings)

			if tt.wantMsg == "" {
				if len(mockNotifier.Messages) > 0 {
					t.Errorf("expected no messages, got %v", mockNotifier.Messages)
				}
			} else {
				if len(mockNotifier.Messages) == 0 {
					t.Errorf("expected message containing %q, got none", tt.wantMsg)
				} else {
					got := mockNotifier.Messages[0]
					if !strings.Contains(got, tt.wantMsg) {
						t.Errorf("expected message containing %q, got %q", tt.wantMsg, got)
					}
					// Check for multiple if needed
					if tt.name == "Multiple Warnings" {
						if !strings.Contains(got, "Battery Voltage Low") {
							t.Errorf("expected message containing 'Battery Voltage Low', got %q", got)
						}
						if !strings.Contains(got, "High Heat Sink Temperature: 90") {
							t.Errorf("expected message containing 'High Heat Sink Temperature: 90', got %q", got)
						}
					}
				}
			}
		})
	}
}

func TestWarningMonitor_ModeChange(t *testing.T) {
	mockNotifier := &MockNotifier{}
	mockStore := NewMockStore()
	wm := NewWarningMonitor(mockNotifier, mockStore)

	// 1. First run, unknown mode -> sets mode, no notification
	pigs := &voltronic.DeviceGeneralStatus{}
	warnings := &voltronic.DeviceWarning{}

	wm.Check(pigs, "line_mode", warnings)
	if len(mockNotifier.Messages) != 0 {
		t.Errorf("expected no notification on initial mode set, got %v", mockNotifier.Messages)
	}
	if val, _ := mockStore.Data["mode"]; val != "line_mode" {
		t.Errorf("expected mode to be saved as line_mode, got %v", val)
	}

	// 2. Second run, same mode -> no notification
	wm.Check(pigs, "line_mode", warnings)
	if len(mockNotifier.Messages) != 0 {
		t.Errorf("expected no notification on same mode, got %v", mockNotifier.Messages)
	}

	// 3. Third run, different mode -> notification
	wm.Check(pigs, "battery_mode", warnings)
	if len(mockNotifier.Messages) != 1 {
		t.Errorf("expected 1 notification on mode change, got %d", len(mockNotifier.Messages))
	} else {
		expected := "Mode changed from line_mode to battery_mode"
		if mockNotifier.Messages[0] != expected {
			t.Errorf("expected notification %q, got %q", expected, mockNotifier.Messages[0])
		}
	}
	if val, _ := mockStore.Data["mode"]; val != "battery_mode" {
		t.Errorf("expected mode to be updated to battery_mode, got %v", val)
	}
}

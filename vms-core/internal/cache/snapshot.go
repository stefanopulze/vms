package cache

import (
	"sync"
	"vms-core/internal/voltronic"
)

type QuerySnapshot struct {
	mux           sync.RWMutex
	ratingInfo    *voltronic.DeviceRatingInfo
	generalStatus *voltronic.DeviceGeneralStatus
	warnings      *voltronic.DeviceWarning
	mode          string
}

func NewQuerySnapshot() *QuerySnapshot {
	return &QuerySnapshot{}
}

func (qs *QuerySnapshot) SetRatingInfo(ratingInfo *voltronic.DeviceRatingInfo) {
	qs.mux.Lock()
	qs.ratingInfo = ratingInfo
	qs.mux.Unlock()
}

func (qs *QuerySnapshot) SetGeneralStatus(generalStatus *voltronic.DeviceGeneralStatus) {
	qs.mux.Lock()
	qs.generalStatus = generalStatus
	qs.mux.Unlock()
}

func (qs *QuerySnapshot) SetWarnings(warnings *voltronic.DeviceWarning) {
	qs.mux.Lock()
	qs.warnings = warnings
	qs.mux.Unlock()
}

func (qs *QuerySnapshot) GetRatingInfo() *voltronic.DeviceRatingInfo {
	qs.mux.RLock()
	defer qs.mux.RUnlock()
	return qs.ratingInfo
}

func (qs *QuerySnapshot) GetGeneralStatus() *voltronic.DeviceGeneralStatus {
	qs.mux.RLock()
	defer qs.mux.RUnlock()
	return qs.generalStatus
}

func (qs *QuerySnapshot) GetWarnings() *voltronic.DeviceWarning {
	qs.mux.RLock()
	defer qs.mux.RUnlock()
	return qs.warnings
}

func (qs *QuerySnapshot) SetMode(mode string) {
	qs.mux.Lock()
	qs.mode = mode
	qs.mux.Unlock()
}

func (qs *QuerySnapshot) GetMode() string {
	qs.mux.RLock()
	defer qs.mux.RUnlock()
	return qs.mode
}

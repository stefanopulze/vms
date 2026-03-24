package domain

import (
	"context"
	"time"
)

type DailyUsage struct {
	ID              int64
	Timestamp       time.Time
	Line            int32 `json:"line"`
	Home            int32 `json:"home"`
	Garage          int32 `json:"garage"`
	Inverter        int32 `json:"inverter"`
	SolarProduction int32 `json:"solarProductionProduction"`
	OffGridPercent  int32 `json:"offGridPercent"`
}

type DailyUsageRepository interface {
	Insert(ctx context.Context, dailyUsage DailyUsage) (int, error)
}

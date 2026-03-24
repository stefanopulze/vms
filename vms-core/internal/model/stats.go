package model

import "time"

type DailyUsage struct {
	Timestamp       time.Time `json:"timestamp"`
	Line            int32     `json:"line"`
	Home            int32     `json:"home"`
	Garage          int32     `json:"garage"`
	Inverter        int32     `json:"inverter"`
	SolarProduction int32     `json:"solarProductionProduction"`
	OffGridPercent  int32     `json:"offGridPercent"`
}

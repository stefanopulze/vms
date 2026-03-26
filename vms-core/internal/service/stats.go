package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"vms-core/internal/domain"
	"vms-core/internal/infrastructure/exporter/influx"
	"vms-core/internal/model"
	"vms-core/internal/repository"
	"vms-core/internal/utils"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

func NewStats(ic *influx.Client, dur *repository.DailyUsage) *Stats {
	return &Stats{
		influx:     ic,
		repository: dur,
	}
}

type Stats struct {
	influx     *influx.Client
	repository *repository.DailyUsage
}

func (s Stats) GetDayPowerUsage(ctx context.Context, day string) (*model.DailyUsage, error) {
	timestamp, err := time.Parse("2006-01-02", day)
	if err != nil {
		return nil, err
	}

	home, err := s.getCategoryPowerUsage(ctx, day, "home")
	if err != nil {
		return nil, err
	}

	line, err := s.getCategoryPowerUsage(ctx, day, "enel")
	if err != nil {
		return nil, err
	}

	inverter, solar, err := s.getInverterPowerUsage(ctx, day)
	if err != nil {
		return nil, err
	}

	total := inverter
	offGrid := int32(0)
	garage := int32(0)

	if inverter == 0 {
		// inverter is 0, so it's turned off
		total = line
	} else {
		offGrid = utils.NonZero(solar * 100 / total)
	}

	garage = utils.NonZero(total - home)

	return &model.DailyUsage{
		Timestamp:       timestamp,
		Home:            home,
		Line:            line,
		Inverter:        inverter,
		Garage:          garage,
		SolarProduction: solar,
		OffGridPercent:  offGrid,
	}, nil
}

func (s Stats) getCategoryPowerUsage(ctx context.Context, day, name string) (int32, error) {
	// sum(active_power) / 43_200 (power read) * 24 (hours); for math simplification 43_200 / 24 = 1800
	query := `SELECT date_bin(interval '1 day', time) as time, sum(active_power) / 1800 as 'active_power'
FROM power_meter 
WHERE time >= DATE_TRUNC('day', $day) AND time < DATE_TRUNC('day', $day) + INTERVAL '1 day' AND name = $name
group by 1
order by time`

	data, err := s.influx.Query(ctx, query, map[string]any{"day": day, "name": name})
	if err != nil {
		return 0, err
	}

	usage := int32(0)
	if len(data) > 0 {
		usage = utils.ConvertAndRoundWatt(data[0]["active_power"])
	}

	return usage, nil
}

func (s Stats) getInverterPowerUsage(ctx context.Context, day string) (int32, int32, error) {
	// 86_400 seconds in a day / read every 5 seconds = 17,280 reads per day * 24 hours = 43_200
	query := `SELECT 
			date_bin(interval '1 day', time) as time, 
			sum(ac_output_active_power ) / 720  as 'ac_output_active_power',
			sum(pv_charging_power) / 720 as 'pv_charging_power'
		FROM inverter 
		WHERE time >= DATE_TRUNC('day', $day) AND time < DATE_TRUNC('day', $day) + INTERVAL '1 day'
		group by 1
		order by time`

	data, err := s.influx.Query(ctx, query, map[string]any{"day": day})
	if err != nil {
		return 0, 0, err
	}

	usage := int32(0)
	solar := int32(0)
	if len(data) > 0 {
		usage = utils.ConvertAndRoundWatt(data[0]["ac_output_active_power"])
		solar = utils.ConvertAndRoundWatt(data[0]["pv_charging_power"])
	}

	return usage, solar, nil
}

func (s Stats) WriteDailyUsage(ctx context.Context, stats *model.DailyUsage) error {
	var err error
	//if err = s.writeInfluxDailyUsage(ctx, stats); err != nil {
	//	slog.Error(err.Error())
	//}
	if err = s.writeDatabaseDailyUsages(ctx, stats); err != nil {
		slog.Error(err.Error())
	}

	return err
}

func (s Stats) writeDatabaseDailyUsages(ctx context.Context, stats *model.DailyUsage) error {
	date := stats.Timestamp.Format("2006-01-02")
	count, err := s.repository.CountByDate(ctx, date)
	if err != nil {
		return err
	}

	du := domain.DailyUsage{
		Home:            stats.Home,
		Line:            stats.Line,
		Inverter:        stats.Inverter,
		Garage:          stats.Garage,
		SolarProduction: stats.SolarProduction,
		OffGridPercent:  stats.OffGridPercent,
		Timestamp:       stats.Timestamp,
	}

	if count > 0 {
		err = s.repository.Update(ctx, date, du)
	} else {
		_, err = s.repository.Insert(ctx, du)
	}

	return err
}

func (s Stats) writeInfluxDailyUsage(ctx context.Context, stats *model.DailyUsage) error {
	tags := map[string]string{}
	fields := map[string]any{
		"home":             stats.Home,
		"line":             stats.Line,
		"inverter":         stats.Inverter,
		"garage":           stats.Garage,
		"solar_production": stats.SolarProduction,
		"offgrid_percent":  stats.OffGridPercent,
	}

	point := influxdb3.NewPoint("daily_usage", tags, fields, stats.Timestamp)

	return s.influx.WritePoints(ctx, point)
}

// DownsamplingDay get day power usage and store it in a database
func (s Stats) DownsamplingDay(ctx context.Context, day string) (*model.DailyUsage, error) {
	slog.Info(fmt.Sprintf("Downsampling day: %s", day))
	stats, err := s.GetDayPowerUsage(ctx, day)
	if err != nil {
		return nil, err
	}

	if err = s.WriteDailyUsage(ctx, stats); err != nil {
		return nil, err
	}

	return stats, nil
}

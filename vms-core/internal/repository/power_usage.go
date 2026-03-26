package repository

import (
	"context"
	"vms-core/internal/domain"
	"vms-core/internal/infrastructure/database"
)

var _ domain.DailyUsageRepository = (*DailyUsage)(nil)

func NewDailyUsage(db *database.Postgres) *DailyUsage {
	return &DailyUsage{
		db: db,
	}
}

type DailyUsage struct {
	db *database.Postgres
}

func (d DailyUsage) CountByDate(ctx context.Context, date string) (int, error) {
	query := `SELECT COUNT(*) FROM daily_usage WHERE date = $1`
	var count int
	err := d.db.Conn().QueryRow(ctx, query, date).Scan(&count)
	return count, err
}

func (d DailyUsage) Insert(ctx context.Context, du domain.DailyUsage) (int, error) {
	query := `INSERT INTO daily_usage (date, line, home, garage, inverter, solar_production, offgrid_percent) 
			VALUES ($1, $2, $3, $4, $5, $6, $7) 
			RETURNING id`

	var id int
	err := d.db.Conn().QueryRow(ctx, query, du.Timestamp.Format("2006-01-02"), du.Line, du.Home, du.Garage, du.Inverter, du.SolarProduction, du.OffGridPercent).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d DailyUsage) Update(ctx context.Context, date string, du domain.DailyUsage) error {
	query := `UPDATE daily_usage SET home = $2, garage = $3, inverter = $4, solar_production = $5, offgrid_percent = $6, line = $7 WHERE date = $1`
	_, err := d.db.Conn().Exec(ctx, query, date, du.Home, du.Garage, du.Inverter, du.SolarProduction, du.OffGridPercent, du.Line)
	return err
}

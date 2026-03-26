package service

import (
	"fmt"
	"log/slog"
	"time"
)

func NewDownsampling(enabled bool) *Downsampling {
	return &Downsampling{
		enabled: enabled,
	}
}

type Downsampling struct {
	enabled bool
}

func (d *Downsampling) Tick(job func()) {
	if !d.enabled {
		slog.Info("Downsampling is disabled")
		return
	}

	slog.Info("Starting downsampling at 01:00 every day")

	for {
		waitDuration := durationUntilNext1AM()
		slog.Info(fmt.Sprintf("Next run in %s (at %s)",
			waitDuration.Round(time.Second),
			time.Now().Add(waitDuration).Format("2006-01-02 15:04:05"),
		))

		timer := time.NewTimer(waitDuration)
		<-timer.C

		go job()

		time.Sleep(1 * time.Second)
	}
}

func durationUntilNext1AM() time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, now.Location())
	if !now.Before(next) {
		next = next.Add(24 * time.Hour)
	}
	return time.Until(next)
}

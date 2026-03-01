package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vms-core/internal/api"
	"vms-core/internal/api/handler"
	"vms-core/internal/cache"
	"vms-core/internal/config"
	"vms-core/internal/event"
	"vms-core/internal/infrastructure/exporter"
	"vms-core/internal/infrastructure/exporter/clickhouse"
	"vms-core/internal/infrastructure/exporter/influx"
	"vms-core/internal/infrastructure/http"
	"vms-core/internal/notifier"
	"vms-core/internal/scheduler"
	"vms-core/internal/serial"
	"vms-core/internal/service"
	"vms-core/internal/voltronic"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	slog.Info("VMS-core")

	//port := testutils.NewDummySerial()
	//testutils.MockStandardCommands(port)
	port, err := serial.NewQueue(&serial.QueueOptions{
		PortName:     cfg.Serial.PortName,
		PortBaudRate: cfg.Serial.BaudRate,
		Size:         cfg.Serial.QueueSize,
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	port.Start()

	inverter := voltronic.NewClient(port)

	// notifier
	tn := notifier.NewTelegram(cfg.Telegram)
	notify := notifier.NewNotify(tn)

	// event manager
	em := event.NewManager()

	// exporters
	exps := exporter.NewMultiple()
	if cfg.Influx.Enabled {
		influxExporter, err := influx.NewClient(cfg.Influx)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		exps.AddExporter(influxExporter)
	}

	if cfg.ClickHouse.Enabled {
		clickhouseExporter, err := clickhouse.NewClient(cfg.ClickHouse)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		exps.AddExporter(clickhouseExporter)
	}

	qs := cache.NewQuerySnapshot()

	// service
	ex := service.NewExporter(inverter, exps, em, qs)

	server := http.NewServer()
	inverterHandler := handler.NewInverter(inverter, em, qs)
	api.BindApi(cfg.Server, server.Router(), port, inverterHandler)
	server.Start()

	// scheduler
	sh := scheduler.NewScheduler(5 * time.Second)
	sh.Tick(ex.ReadStatusInformation)
	sh.Start()

	_ = notify.Send(ctx, "VMS-core started")

	// Listen for the interrupt signal
	<-ctx.Done()
	slog.Info("Shutting down...")
	stop()
	sh.Stop()
	server.Stop()
	_ = exps.Close()
	_ = port.Close()
}

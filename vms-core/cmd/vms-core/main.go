package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
	"vms-core/internal/api"
	"vms-core/internal/api/handler"
	"vms-core/internal/cache"
	"vms-core/internal/config"
	"vms-core/internal/infrastructure/database"
	"vms-core/internal/infrastructure/exporter/influx"
	"vms-core/internal/infrastructure/http"
	"vms-core/internal/infrastructure/telegram"
	"vms-core/internal/notifier"
	"vms-core/internal/repository"
	"vms-core/internal/scheduler"
	"vms-core/internal/serial"
	"vms-core/internal/service"
	"vms-core/internal/service/commands"
	"vms-core/internal/store"
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

	tc := telegram.NewClient(cfg.Telegram)

	// notifier
	notify := notifier.NewNotify(notifier.NewTelegram(tc))

	// exporters
	influxClient, err := influx.NewClient(cfg.Influx)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	qs := cache.NewQuerySnapshot()

	if cfg.Telegram.EnableCommands {
		statusCommand := commands.NewStatusCommand(tc, qs)
		updateSourceCommand := commands.NewUpdateSourcePriority(tc, inverter)
		remoteCommandService := service.NewRemoteCommands(tc, statusCommand, updateSourceCommand)
		go tc.GetUpdates(ctx, remoteCommandService.HandleTelegramCommand)
	}

	// database
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		slog.Error("cannot create postgres client", slog.Any("error", err))
		os.Exit(1)
	}

	// repository
	dur := repository.NewDailyUsage(db)

	// service
	storePath := path.Join(cfg.Storage, "vms_state.json")
	storage, err := store.NewFileStore(storePath)
	if err != nil {
		slog.Error("cannot create file store", slog.Any("error", err))
		os.Exit(1)
	}

	wm := service.NewWarningMonitor(notify, storage)
	scheduledCommands := service.NewScheduledCommands(inverter, influxClient, qs, wm)

	statsService := service.NewStats(influxClient, dur)

	server := http.NewServer()
	inverterHandler := handler.NewInverter(inverter, qs)
	statsHandler := handler.NewStats(influxClient, statsService)
	api.BindApi(cfg.Server, server.Router(), port, inverterHandler, statsHandler)
	server.Start()

	// scheduler
	sh := scheduler.NewScheduler(5 * time.Second)
	sh.Tick(scheduledCommands.Read)
	sh.Start()

	if cfg.Downsampling.Enabled {
		downsampleService := service.NewDownsampling(cfg.Downsampling.Enabled)
		downsampleService.Tick(func() {
			yesterday := time.Now().AddDate(0, 0, -1)
			if _, statsErr := statsService.DownsamplingDay(ctx, yesterday.Format("2006-01-02")); statsErr != nil {
				slog.Error("cannot downsample yesterday", slog.Any("error", err))
			}
		})
	}

	_ = notify.Send(ctx, "VMS-core started")

	// Listen for the interrupt signal
	<-ctx.Done()
	slog.Info("Shutting down...")
	stop()
	sh.Stop()
	server.Stop()
	_ = influxClient.Close()
	_ = port.Close()
	_ = db.Close()
}

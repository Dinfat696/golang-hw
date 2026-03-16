package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/config"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/kafka"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/metrics" // <-- ДОБАВЛЕНО
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
	memorystorage "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage/memory" // <-- РАСКОММЕНТИРОВАНО
	sqlstorage "github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logg, err := logger.NewLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	var store storage.Storage
	if cfg.Storage.Type == "sql" {
		// Используем sqlstorage (алиас для internal/storage/sql)
		store, err = sqlstorage.NewStorage(cfg.Storage.DSN) // <-- ИСПРАВЛЕНО: sql. → sqlstorage.
		if err != nil {
			logg.Fatalf("Failed to create SQL storage: %v", err)
		}
	} else {
		// Используем memorystorage (алиас для internal/storage/memory)
		store = memorystorage.NewStorage() // <-- ИСПРАВЛЕНО: memory. → memorystorage.
	}
	defer func() {
    if err := store.Close(); err != nil {
        logg.Errorf("failed to close storage: %v", err)
    }
}()

	// Создаем и подключаем Kafka producer с retry
	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := producer.WaitForConnect(ctx, cfg.Kafka.MaxAttempts, cfg.Kafka.RetryBackoff); err != nil {
		logg.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer func() {
    if err := store.Close(); err != nil {
        logg.Errorf("failed to close storage: %v", err)
    }
}()

	logg.Info("Successfully connected to Kafka")

	calendarApp := app.New(logg, store)
	metricsInstance := metrics.NewMetrics() // <-- ТЕПЕРЬ РАБОТАЕТ, ТАК КАК ДОБАВЛЕН ИМПОРТ
	scheduler := app.NewScheduler(calendarApp, producer, logg, cfg.Scheduler, metricsInstance)

	// Graceful shutdown
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logg.Info("Shutting down scheduler...")
		mainCancel()
	}()

	logg.Info("Starting scheduler...")
	if err := scheduler.Run(mainCtx); err != nil {
		logg.Errorf("Scheduler error: %v", err)
		os.Exit(1)
	}
}
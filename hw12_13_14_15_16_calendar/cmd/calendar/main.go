package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/api" // для api.Handler
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/server/http"    // если internalhttp находится здесь
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage/memory" // для memorystorage
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage/sql"    // для sqlstorage
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config",
		"D:\\GolandProjects\\golang-hw\\hw12_13_14_15_calendar\\configs\\config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Level, config.Location)
	var calendar *app.App
	if config.Storage == "inner" {
		storage := memorystorage.New()
		calendar = app.New(logg, storage)
	} else {
		storage := sqlstorage.New(config.DbDriverName, config.Dsn)
		calendar = app.New(logg, storage)
	}
	fmt.Print(calendar)

	handler := internalhttp.NewEventsHandler(logg, calendar, config.Host, config.Port)

	httpHandler := api.Handler(handler)

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Port),
		Handler:      httpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Сервер запущен на :8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

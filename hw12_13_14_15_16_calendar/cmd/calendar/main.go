package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

	// Создаём экземпляр хэндлера
	handler := internalhttp.NewEventsHandler(logg, calendar, config.Host, config.Port)

	// Генерируем HTTP-хэндлер из oapi-codegen
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
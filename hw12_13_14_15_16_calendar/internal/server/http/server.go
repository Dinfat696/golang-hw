package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/logger"
)

type Server struct {
	host        string
	port        int
	logger      *logger.Logger
	application *app.App
}

func NewServer(logger *logger.Logger, app *app.App, host string, port int) *Server {
	return &Server{
		host:        host,
		port:        port,
		logger:      logger,
		application: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/openapi.yaml", swagger)

	srv := &http.Server{
		Addr:         s.host + ":" + strconv.Itoa(s.port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,  
		WriteTimeout: 10 * time.Second, 
		IdleTimeout:  60 * time.Second, 
		
	}

	
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		return nil

	case err := <-errCh:
		
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed to start or crashed: %w", err)
		}
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	fmt.Println(ctx)
	return nil
}

func swagger(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "openapi.yaml")
}

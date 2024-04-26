package main

import (
	"github.com/go-chi/httplog/v2"
	"log/slog"
	"sync"
	"time"
)

const (
	version = "1.0.0"
)

type application struct {
	config *config
	logger *httplog.Logger
	wg     sync.WaitGroup
}

func main() {
	cfg := newConfig()

	logger := httplog.NewLogger("Guild-Hall-Server", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		Tags: map[string]string{
			"version": version,
			"env":     cfg.env,
		},
		QuietDownRoutes: []string{
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
	})

	app := &application{
		config: cfg,
		logger: logger,
	}
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		app.logger.Info("Hello from the go function!")
		app.logger.Debug("Oh no, a debug message!")
	}()

	app.wg.Wait()
}

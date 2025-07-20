package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meh-hackathon/meh/config"
	"github.com/meh-hackathon/meh/db"
	"github.com/meh-hackathon/meh/logger"
	"github.com/meh-hackathon/meh/server"
)

func main() {
	start := time.Now()

	err := config.Load()
	mustSetup(err)

	err = setupLogger(config.LogLevel)
	mustSetup(err)

	err = db.Init(config.DatabaseURL)
	mustSetup(err)

	srv, err := server.New()
	mustSetup(err)

	logger.Info("setup done", "after", time.Since(start))

	srv.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	logger.Debug("Shutting down application...", "uptime", time.Since(start))

	err = srv.Shutdown()
	mustShutdown(err)

	err = db.Close()
	mustShutdown(err)
}

func mustSetup(err error) {
	if err != nil {
		panic(fmt.Errorf("Fatal setup error: %w", err))
	}
}

func mustShutdown(err error) {
	if err != nil {
		logger.Error("Fatal shutdown error", "error", err)
		os.Exit(1)
	}
}

func setupLogger(level slog.Level) error {
	start := time.Now()
	filename := "logs/meh.log.jsonl"

	logger.New()
	logger.SetLevel(level)

	logger.Debug("Logger initialized", "filename", filename, "level", level, "took", time.Since(start))
	return nil
}

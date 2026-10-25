package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JesusEugenio/MitSys/internal/config"
	"github.com/JesusEugenio/MitSys/internal/db"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("error executing script", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to loading configuration: %w", err)
	}

	startupCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	pool, err := db.Start(startupCtx, cfg.DBUrl())
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		slog.Info("closing database connections")
		pool.Close()
	}()

	slog.Info("connected to MitSys db")

	return nil
}

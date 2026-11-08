package main

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JesusEugenio/MitSys/internal/config"
	"github.com/JesusEugenio/MitSys/internal/db"
	"github.com/JesusEugenio/MitSys/internal/handlers"
	"github.com/JesusEugenio/MitSys/internal/router"
	"github.com/JesusEugenio/MitSys/internal/service"
	"github.com/JesusEugenio/MitSys/internal/store"
	"github.com/JesusEugenio/MitSys/internal/tmpls"
	"github.com/JesusEugenio/MitSys/ui"
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

	tmplsFS, err := fs.Sub(ui.Embed, "tmpls")
	if err != nil {
		return fmt.Errorf("error accesing tmpl FS: %w", err)
	}

	tmplCache, err := tmpls.LoadCache(tmplsFS)
	if err != nil {
		return fmt.Errorf("error loading tmpls: %w", err)
	}

	InventoryStore := store.NewInventory(pool)
	InventoryService := service.NewInvetoryService(InventoryStore)
	InventoryHandler := handlers.NewInventoryHandler(InventoryService, tmplCache)

	routerServer := router.NewServer(InventoryHandler)
	r := routerServer.Setup()

	address := ":" + cfg.Port
	slog.Info(fmt.Sprintf("MitSys running on http://localhost%s", address))

	server := &http.Server{
		Addr:    address,
		Handler: r,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("critic fail on http server: %w", err)
		}
		return nil

	case <-ctx.Done():
		slog.Info("shutdowing http server ...")

		shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdownCtx()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("error shutdowing server: %w", err)
		}

		slog.Info("server shutdown correctly")
		return nil
	}
}

package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/JesusEugenio/MitSys/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
)

func Start(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating pgxpool: %w", err)
	}
	defer func() {
		if err != nil {
			pool.Close()
		}
	}()

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging postgres pool: %w", err)
	}

	if err = migrate(ctx, pool); err != nil {
		return nil, err
	}

	slog.Info("applied migrations, db connected with pgxpool")
	return pool, nil
}

func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		if err := db.Close(); err != nil {
			slog.Warn("error closing db wrapper", "error", err)
		}
	}()

	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return fmt.Errorf("error creating postgres session locker: %w", err)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		migrations.Embed,
		goose.WithSessionLocker(locker),
		goose.WithSlog(slog.Default()),
	)
	if err != nil {
		return fmt.Errorf("error creating goose provider: %w", err)
	}

	slog.Info("executing db migrations with goose")
	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("error executing migrations with goose: %w", err)
	}

	return nil
}

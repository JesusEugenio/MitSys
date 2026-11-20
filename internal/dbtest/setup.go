package dbtest

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func Setup(ctx context.Context) (*pgxpool.Pool, func() error, error) {
	dbName := "mitsys_test_db"
	dbUser := "tester"
	dbPassword := "test123"

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error running pgContainer: %w", err)
	}

	terminate := func() error {
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			return fmt.Errorf("error in terminate container for pgContainer: %w", err)
		}
		return nil
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		err = fmt.Errorf("error getting connection string of pgContainer: %w", err)
		return nil, nil, errors.Join(err, terminate())
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		err = fmt.Errorf("error creating new pgxpool of pgContainer: %w", err)
		return nil, nil, errors.Join(err, terminate())
	}

	cleanup := func() error {
		pool.Close()
		return terminate()
	}

	return pool, cleanup, nil
}

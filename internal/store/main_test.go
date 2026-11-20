package store

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/JesusEugenio/MitSys/internal/db"
	"github.com/JesusEugenio/MitSys/internal/dbtest"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var (
		cleanup func() error
		err     error
	)

	testPool, cleanup, err = dbtest.Setup(context.Background())
	if err != nil {
		slog.Error("error in setup db", "error", err)
		os.Exit(1)
	}

	err = db.Migrate(context.Background(), testPool)
	if err != nil {
		slog.Error("error doing pool migrate", "error", err)
		os.Exit(1)
	}

	code := m.Run()

	if err = cleanup(); err != nil {
		slog.Warn("error cleaning test pool and container", "warning", err)
	}

	os.Exit(code)
}

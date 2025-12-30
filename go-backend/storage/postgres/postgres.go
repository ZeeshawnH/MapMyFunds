package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	// Disable prepared statement cache to avoid conflicts with Air hot-reload
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	return pgx.ConnectConfig(ctx, config)
}

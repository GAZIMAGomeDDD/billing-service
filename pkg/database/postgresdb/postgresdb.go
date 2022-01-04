package postgresdb

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB(ctx context.Context, connString string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres due to error: %w", err)
	}

	return db, nil
}

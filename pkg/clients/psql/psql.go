package psql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Addr     string
	Port     string
	Username string
	Password string
	Database string
}

func NewClient(ctx context.Context, cfg PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Addr,
		cfg.Port,
		cfg.Database,
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

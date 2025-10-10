package postgres

import (
	"context"
	"fmt"
	"go-data-catalog/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(cfg *config.Config) (*DB, error) {
	connString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable&client_encoding=utf8",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
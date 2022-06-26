package postgres

import (
	"context"
	"fmt"

	"github.com/MarySmirnova/comments_service/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

var ctx = context.Background()

type Store struct {
	db *pgxpool.Pool
}

func New(cfg config.Postgres) (*Store, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) GetPGXPool() *pgxpool.Pool {
	return s.db
}

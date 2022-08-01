package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"gophermart/internal/domain/models"
	"gophermart/internal/ports"
)

type PostgresStore struct {
	db     *pgxpool.Pool
	logger ports.Logger
}

func NewPostgresStore(ctx context.Context, pgconn string, logger ports.Logger) (*PostgresStore, error) {
	var store PostgresStore

	store.logger = logger

	config, err := pgxpool.ParseConfig(pgconn)
	if err != nil {
		store.logger.Error("postgres connection config parse failed", err)
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		store.logger.Error("postgres pool creation failed", err)
		return nil, err
	}

	store.db = pool

	return &store, nil
}

func (s PostgresStore) RegisterUser(u models.User) error {
	s.logger.Info("register user:", u.Login)
	return nil
}

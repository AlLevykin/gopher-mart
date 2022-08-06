package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gophermart/internal/domain/models"
	"gophermart/internal/domain/repo"
	"gophermart/internal/ports"
)

//go:embed pg-migrations/*.sql
var embedMigrations embed.FS

type PostgresStore struct {
	db     *sql.DB
	logger ports.Logger
}

func NewPostgresStore(ctx context.Context, pgconn string, logger ports.Logger) (*PostgresStore, error) {
	var store PostgresStore

	store.logger = logger

	goose.SetBaseFS(embedMigrations)
	db, err := goose.OpenDBWithDriver("postgres", pgconn)
	if err != nil {
		store.logger.Error("postgres pool creation failed", err)
		return nil, err
	}
	if err := goose.Up(db, "pg-migrations"); err != nil {
		store.logger.Error("migrations failed", err)
		db.Close()
		return nil, err
	}

	store.db = db

	return &store, nil
}

func (s PostgresStore) RegisterUser(ctx context.Context, u *models.User) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO \"user\"(login, pwh, salt) VALUES($1,$2,$3)",
		u.Login, u.Password, "")

	if err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == "23505" {
			s.logger.Error("user registration filed:", pgerr.Detail)
			return repo.ErrUserExists
		}
		s.logger.Error("user registration filed:", err)
		return err
	}

	s.logger.Info("register user:", u.Login)
	return nil
}

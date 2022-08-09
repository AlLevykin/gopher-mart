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

func (s PostgresStore) Validation(ctx context.Context, u *models.User) error {
	s.logger.Info("user validation:", u.Login)
	r, err := s.db.ExecContext(
		ctx,
		"SELECT * FROM \"user\" WHERE login = $1 AND pwh = $2",
		u.Login,
		u.Password)
	if err != nil {
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return repo.ErrUserValidation
	}
	return nil
}

func (s PostgresStore) IsOrderAccepted(ctx context.Context, order string) (bool, error) {
	s.logger.Info("order checking:", order)
	r, err := s.db.ExecContext(
		ctx,
		"SELECT * FROM \"order\" WHERE \"number\" = $1",
		order)
	if err != nil {
		return false, err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows != 0 {
		return false, nil
	}
	return true, nil
}

func (s PostgresStore) IsOrderAcceptedByUser(ctx context.Context, order string, login string) (bool, error) {
	s.logger.Info("order checking by user:", order, " ", login)
	r, err := s.db.ExecContext(
		ctx,
		"SELECT * FROM \"order\" WHERE \"number\" = $1 AND \"user\" = $2",
		order,
		login)
	if err != nil {
		return false, err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows != 0 {
		return false, nil
	}
	return true, nil
}

func (s PostgresStore) SaveOrder(ctx context.Context, order string, login string) error {
	s.logger.Info("save order:", order, " ", login)
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO \"order\"(\"number\", \"user\", status) VALUES($1,$2,$3)",
		order, login, "NEW")
	if err != nil {
		s.logger.Error("can't save order:", err)
		return err
	}
	return nil
}

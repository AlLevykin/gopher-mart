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
		"INSERT INTO \"order\"(\"number\", \"user\", status, accrual) VALUES($1,$2,$3,$4)",
		order, login, "PROCESSED", 729.98)
	if err != nil {
		s.logger.Error("can't save order:", err)
		return err
	}
	return nil
}

func (s PostgresStore) GetOrders(ctx context.Context, login string) (string, error) {
	var json string
	row := s.db.QueryRowContext(ctx, "SELECT json_agg(row_to_json(row)) AS JSON FROM (SELECT \"number\", status, accrual, to_char(uploaded, 'YYYY-MM-DD\"T\"HH24:MI:SS.US\"Z\"') AS uploaded_at FROM \"order\" WHERE \"user\"=$1) row", login)
	err := row.Scan(&json)
	if err != nil {
		s.logger.Error("can't get orders:", err)
		return "", sql.ErrNoRows
	}
	return json, nil
}

func (s PostgresStore) GetBalance(ctx context.Context, login string) (string, error) {
	var json string
	row := s.db.QueryRowContext(ctx, "SELECT row_to_json(row) AS JSON FROM (SELECT \"current\", \"withdrawn\" FROM \"balance\" WHERE \"user\"=$1) row", login)
	err := row.Scan(&json)
	if err != nil {
		s.logger.Error("can't get orders:", err)
		return "", sql.ErrNoRows
	}
	return json, nil
}

func (s PostgresStore) SaveWithdraw(ctx context.Context, w *models.Withdraw) error {
	s.logger.Info("save withdraw:", w.Order, " ", w.Sum)

	_, err := s.db.ExecContext(ctx,
		"INSERT INTO \"withdraw\"(\"order\", \"sum\") VALUES($1,$2)",
		w.Order, w.Sum)
	if err != nil {
		s.logger.Error("can't save withdraw:", err)
		return err
	}
	return nil
}

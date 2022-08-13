package ports

import (
	"context"
	"gophermart/internal/domain/models"
)

type Store interface {
	RegisterUser(ctx context.Context, u *models.User) error
	Validation(ctx context.Context, u *models.User) error
	IsOrderAccepted(ctx context.Context, order string) (bool, error)
	IsOrderAcceptedByUser(ctx context.Context, order string, login string) (bool, error)
	SaveOrder(ctx context.Context, order string, login string) error
	GetOrders(ctx context.Context, login string) (string, error)
	GetBalance(ctx context.Context, login string) (string, error)
	SaveWithdraw(ctx context.Context, w *models.Withdraw) error
	GetWithdrawals(ctx context.Context, login string) (string, error)
}

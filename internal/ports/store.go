package ports

import (
	"context"
	"gophermart/internal/domain/models"
)

type Store interface {
	RegisterUser(ctx context.Context, u *models.User) error
	Validation(ctx context.Context, u *models.User) error
}

package ports

import "gophermart/internal/domain/models"

type Store interface {
	RegisterUser(u models.User) error
}

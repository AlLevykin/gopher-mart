package accrualdispatcher

import (
	"context"
	"gophermart/internal/domain/models"
	"gophermart/internal/ports"
)

type GopherAccrualDispatcher struct {
	address string
	store   ports.Store
	logger  ports.Logger
}

func NewGopherAccrualDispatcher(address string, store ports.Store, logger ports.Logger) *GopherAccrualDispatcher {
	logger.Info("accrual dispatcher started")
	return &GopherAccrualDispatcher{
		address: address,
		store:   store,
		logger:  logger,
	}
}

func (d GopherAccrualDispatcher) Dispatch(order string) {
	d.logger.Info("dispatch order:", order)
	d.store.UpdateOrder(context.Background(), &models.Order{Number: order})
}

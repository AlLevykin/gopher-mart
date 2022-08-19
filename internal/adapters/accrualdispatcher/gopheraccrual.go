package accrualdispatcher

import (
	"context"
	"gophermart/internal/domain/repo"
	"gophermart/internal/ports"
	"io/ioutil"
	"net/http"
)

type GopherAccrualDispatcher struct {
	address string
	chReq   chan string
	store   ports.Store
	logger  ports.Logger
}

func NewGopherAccrualDispatcher(address string, store ports.Store, logger ports.Logger) *GopherAccrualDispatcher {
	logger.Info("accrual dispatcher started")
	d := &GopherAccrualDispatcher{
		address: address,
		store:   store,
		logger:  logger,
		chReq:   make(chan string),
	}
	return d
}

func (d GopherAccrualDispatcher) Start() {
	go func() {
		for num := range d.chReq {
			go func(order string) {
				resp, err := http.Get(d.address + "/api/orders/" + order)
				if err != nil {
					d.logger.Error("accrual service error:", err)
					return
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					d.logger.Error("accrual service error:", err)
					return
				}
				if resp.StatusCode != http.StatusOK {
					d.logger.Error("accrual service response status:", resp.StatusCode)
					return
				}
				d.logger.Info(string(body))
				o, err := repo.UnmarshalOrder(string(body))
				if err != nil {
					d.logger.Error("accrual service error:", err)
					return
				}
				d.store.UpdateOrder(context.Background(), o)
			}(num)
		}
	}()
}

func (d GopherAccrualDispatcher) Stop() {
	close(d.chReq)
	d.logger.Info("accrual dispatcher stopped")
}

func (d GopherAccrualDispatcher) Dispatch(order string) {
	d.chReq <- order
}

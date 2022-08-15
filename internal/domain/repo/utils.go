package repo

import (
	"encoding/json"
	"gophermart/internal/domain/models"
)

func UnmarshalUser(s string) (*models.User, error) {
	res := models.User{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func UnmarshalWithdraw(s string) (*models.Withdraw, error) {
	res := models.Withdraw{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func UnmarshalBalance(s string) (*models.Balance, error) {
	res := models.Balance{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func UnmarshalOrder(s string) (*models.Order, error) {
	res := models.Order{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

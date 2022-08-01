package repo

import (
	"encoding/json"
	"gophermart/internal/domain/models"
)

func UnmarshalUser(s string) (models.User, error) {
	res := models.User{}
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		return models.User{}, err
	}
	return res, nil
}

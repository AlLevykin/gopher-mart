package models

type Balance struct {
	Current   float64 `json:"current" type:"number"`
	Withdrawn float64 `json:"withdrawn" type:"number"`
}

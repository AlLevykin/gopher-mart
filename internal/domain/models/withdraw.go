package models

type Withdraw struct {
	Order       string  `json:"order" type:"string"`
	Sum         float64 `json:"sum" type:"number"`
	ProcessedAt string  `json:"processed_at" type:"string"`
}

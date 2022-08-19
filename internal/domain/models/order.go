package models

type Order struct {
	Number     string  `json:"order" type:"string"`
	Status     string  `json:"status" type:"string"`
	Accrual    float64 `json:"accrual" type:"number"`
	UploadedAt string  `json:"uploaded_at" type:"string"`
}

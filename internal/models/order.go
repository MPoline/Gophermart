package models

import "time"

type Order struct {
	Number     string `json:"number"`
	Login      string
	Status     string    `json:"status"` // NEW, PROCESSING, INVALID, PROCESSED
	Accrual    *float64  `json:"accrual,omitempty"`
	Amount     *float64  `json:"amount,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

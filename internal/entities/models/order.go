package models

import (
	"time"
)

type OrderDTO struct {
	UserID     int       `json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

var (
	StatusNew         = `NEW`
	StatusRegistering = `REGISTERED`
	StatusInvalid     = `INVALID`
	StatusProcessing  = `PROCESSING`
	StatusProcessed   = `PROCESSED`
)

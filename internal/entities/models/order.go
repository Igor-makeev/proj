package models

import (
	"time"
)

type OrderDTO struct {
	UserID     int       `json:"-"`
	Number     string    `json:"number,string"`
	Status     string    `json:"status,string"`
	Accrual    float64   `json:"accrual,string,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
}

var (
	StatusNew         = "NEW"
	StatusRegistering = "REGISTERED"
	StatusInvalid     = "INVALID"
	StatusProcessing  = "PROCESSING"
	StatusProcessed   = "PROCESSED"
)

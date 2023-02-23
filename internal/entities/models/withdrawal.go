package models

import "time"

type Withdrawal struct {
	OrderNumber string    `json:"order" binding:"required"`
	Sum         float64   `json:"sum" binding:"required"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}

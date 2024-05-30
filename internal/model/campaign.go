package model

import (
	"time"

	"github.com/google/uuid"
)

type Campaign struct {
	Id         uuid.UUID `json:"id"`
	MerchantId uuid.UUID `json:"merchant_id"`
	Active     bool      `json:"active"`
	Lat        float64   `json:"lat"`
	Long       float64   `json:"long"`
	CreatedBy  string    `json:"created_by"`
	UpdatedBy  string    `json:"updated_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

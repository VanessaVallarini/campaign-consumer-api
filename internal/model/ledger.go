package model

import (
	"time"

	"github.com/google/uuid"
)

type Ledger struct {
	Id         uuid.UUID `json:"id"`
	CampaignId uuid.UUID `json:"campaign_id"`
	SpentId    uuid.UUID `json:"spent_id"`
	UserId     uuid.UUID `json:"user_id"`
	EventType  string    `json:"event_type"`
	Cost       float64   `json:"cost"`
	Lat        float64   `json:"lat"`
	Long       float64   `json:"long"`
	CreatedAt  time.Time `json:"created_at"`
}

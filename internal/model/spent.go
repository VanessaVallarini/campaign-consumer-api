package model

import "github.com/google/uuid"

type Spent struct {
	Id               uuid.UUID `json:"id"`
	CampaignId       uuid.UUID `json:"campaign_id"`
	Bucket           string    `json:"bucket"`
	TotalSpent       float64   `json:"total_spent"`
	TotalClicks      float64   `json:"total_clicks"`
	TotalImpressions float64   `json:"total_impressions"`
}

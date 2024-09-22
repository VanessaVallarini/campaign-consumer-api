package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	CampaignCreatedAndActive = "Criação de campanha"
	CampaignUpdateStatus     = "Alteração do status de %s para %s"
	CampaignUpdateBudget     = "Alteração do limite diário de R$%.0f para R$%.0f"
)

type CampaignHistory struct {
	Id          uuid.UUID `json:"id"`
	CampaignId  uuid.UUID `json:"campaign_id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

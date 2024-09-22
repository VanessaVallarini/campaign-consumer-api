package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	RegionCreatedAndActive = "Criação de região"
	RegionUpdateStatus     = "Alteração do status de %s para %s"
	RegionUpdateCost       = "Alteração do custo de R$%.0f para R$%.0f"
)

type RegionHistory struct {
	Id          uuid.UUID `json:"id"`
	RegionId    uuid.UUID `json:"region_id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

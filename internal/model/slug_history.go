package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	SlugCreatedAndActive = "Criação de slug"
	SlugUpdateStatus     = "Alteração do status de %s para %s"
	SlugUpdateCost       = "Alteração do limite diário de R$%.0f para R$%.0f"
)

type SlugHistory struct {
	Id          uuid.UUID `json:"id"`
	SlugId      uuid.UUID `json:"slug_id"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

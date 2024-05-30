package model

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	Id        uuid.UUID   `json:"id"`
	OwnerId   uuid.UUID   `json:"owner_id"`
	Slugs     []uuid.UUID `json:"slugs"`
	Name      string      `json:"name"`
	Active    bool        `json:"active"`
	CreatedBy string      `json:"created_by"`
	UpdatedBy string      `json:"updated_by"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

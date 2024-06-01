package model

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	Id        uuid.UUID   `json:"id"`
	OwnerId   uuid.UUID   `json:"owner_id"`
	RegionId  uuid.UUID   `json:"region_id"`
	Slugs     []uuid.UUID `json:"slugs"`
	Name      string      `json:"name"`
	Active    bool        `json:"active"`
	CreatedBy string      `json:"created_by"`
	UpdatedBy string      `json:"updated_by"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type MerchantEvent struct {
	Id        string   `json:"id"`
	OwnerId   string   `json:"owner_id"`
	RegionId  string   `json:"region_id"`
	Slugs     []string `json:"slugs"`
	Name      string   `json:"name"`
	Active    string   `json:"active"` //conseguimos não usar apenas string?
	CreatedBy string   `json:"created_by"`
	UpdatedBy string   `json:"updated_by"`
}

const (
	MerchantAvro = `{
		"type":"record",
		"name":"merchant",
		"namespace":"campaign.campaign_merchant_value",
		"fields":[
			 {
				"name":"id",
				"type":"string"
			 },
			 {
				"name":"owner_id",
				"type":"string"
			 },
			 {
				"name":"region_id",
				"type":"string"
			 },
			 {
				"name":"slugs",
				"type": {
					"type": "array",
					"items": "string"
				}
			 },
			 {
				"name":"name",
				"type":"string"
			 },
			 {
				"name":"active",
				"type":"string"
			 },
			 {
				"name":"created_by",
				"type":"string"
			 },
			 {
				"name":"updated_by",
				"type":"string"
			 }		   
		]
	 }`
)

package model

import (
	"time"

	"github.com/google/uuid"
)

type MerchantStatus string

const (
	ActiveMerchant   MerchantStatus = "ACTIVE"
	InactiveMerchant MerchantStatus = "INACTIVE"
)

type Merchant struct {
	Id        uuid.UUID      `json:"id"`
	OwnerId   uuid.UUID      `json:"owner_id"`
	RegionId  uuid.UUID      `json:"region_id"`
	Slugs     []uuid.UUID    `json:"slugs"`
	Name      string         `json:"name"`
	Status    MerchantStatus `json:"status"`
	CreatedBy string         `json:"created_by"`
	UpdatedBy string         `json:"updated_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type MerchantEvent struct {
	Id        uuid.UUID `avro:"id"`
	OwnerId   uuid.UUID `avro:"owner_id"`
	RegionId  uuid.UUID `avro:"region_id"`
	Slugs     []string  `avro:"slugs"`
	Name      string    `avro:"name"`
	Status    string    `avro:"status"`
	User      string    `avro:"user"`
	EventTime time.Time `avro:"eventTime"`
}

const (
	MerchantAvro = `{
		"type":"record",
		"name":"merchant",
		"namespace":"campaign.campaign_merchant_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "owner_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "region_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
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
				"name":"status",
				"type":"string"
			},
			{
				"name":"user",
				"type":"string"
			},
			{
				"name": "eventTime",
				"type": {
				"type": "long",
				"logicalType": "timestamp-millis"
				}
			}		   
		]
	 }`
)

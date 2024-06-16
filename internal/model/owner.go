package model

import (
	"time"

	"github.com/google/uuid"
)

type OwnerStatus string

const (
	ActiveOwner   OwnerStatus = "ACTIVE"
	InactiveOwner OwnerStatus = "INACTIVE"
)

type Owner struct {
	Id        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	Status    OwnerStatus `json:"status"`
	CreatedBy string      `json:"created_by"`
	UpdatedBy string      `json:"updated_by"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OwnerEvent struct {
	Id        uuid.UUID `avro:"id"`
	Email     string    `avro:"email"`
	Status    string    `avro:"status"`
	User      string    `avro:"user"`
	EventTime time.Time `avro:"even_time"`
}

const (
	OwnerAvro = `{
		"type":"record",
		"name":"owner",
		"namespace":"campaign.campaign_owner_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name":"email",
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
				"name": "even_time",
				"type": {
				"type": "long",
				"logicalType": "timestamp-millis"
				}
			}	   
		]
	 }`
)

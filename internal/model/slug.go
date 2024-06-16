package model

import (
	"time"

	"github.com/google/uuid"
)

type SlugStatus string

const (
	ActiveSlug   SlugStatus = "ACTIVE"
	InactiveSlug SlugStatus = "INACTIVE"
)

type Slug struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Status    SlugStatus `json:"status"`
	Cost      float64    `json:"cost"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy string     `json:"updated_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type SlugEvent struct {
	Id        uuid.UUID `avro:"id"`
	Name      string    `avro:"name"`
	Status    string    `avro:"status"`
	Cost      float64   `avro:"cost"`
	User      string    `avro:"user"`
	EventTime time.Time `avro:"eventTime"`
}

const (
	SlugAvro = `{
		"type":"record",
		"name":"slug",
		"namespace":"campaign.campaign_slug_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
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
				"name":"cost",
				"type":"double"
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

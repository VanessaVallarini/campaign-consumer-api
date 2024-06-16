package model

import (
	"time"

	"github.com/google/uuid"
)

type RegionStatus string

const (
	ActiveRegion   RegionStatus = "ACTIVE"
	InactiveRegion RegionStatus = "INACTIVE"
)

type Region struct {
	Id        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Status    RegionStatus `json:"status"`
	Lat       float64      `json:"lat"`
	Long      float64      `json:"long"`
	Cost      float64      `json:"cost"`
	CreatedBy string       `json:"created_by"`
	UpdatedBy string       `json:"updated_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type RegionEvent struct {
	Id        uuid.UUID `avro:"id"`
	Name      string    `avro:"name"`
	Status    string    `avro:"status"`
	Lat       float64   `avro:"lat"`
	Long      float64   `avro:"long"`
	Cost      float64   `avro:"cost"`
	User      string    `avro:"user"`
	EventTime time.Time `avro:"even_time"`
}

const (
	RegionAvro = `{
		"type":"record",
		"name":"region",
		"namespace":"campaign.campaign_region_value",
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
				"name":"lat",
				"type":"double"
			},
			{
				"name":"long",
				"type":"double"
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
				"name": "even_time",
				"type": {
				"type": "long",
				"logicalType": "timestamp-millis"
				}
			}		   
		]
	 }`
)

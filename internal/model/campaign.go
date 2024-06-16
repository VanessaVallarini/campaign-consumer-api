package model

import (
	"time"

	"github.com/google/uuid"
)

type CampaignStatus string

const (
	ActiveCampaign   CampaignStatus = "ACTIVE"
	InactiveCampaign CampaignStatus = "INACTIVE"
)

type Campaign struct {
	Id         uuid.UUID      `json:"id"`
	MerchantId uuid.UUID      `json:"merchant_id"`
	Status     CampaignStatus `json:"status"`
	Lat        float64        `json:"lat"`
	Long       float64        `json:"long"`
	Budget     float64        `json:"budget"`
	CreatedBy  string         `json:"created_by"`
	UpdatedBy  string         `json:"updated_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type CampaignEvent struct {
	Id         uuid.UUID `avro:"id"`
	MerchantId uuid.UUID `avro:"merchant_id"`
	Status     string    `avro:"status"`
	Lat        float64   `avro:"lat"`
	Long       float64   `avro:"long"`
	Budget     float64   `avro:"budget"`
	User       string    `avro:"user"`
	EventTime  time.Time `avro:"eventTime"`
}

const (
	CampaignAvro = `{
		"type":"record",
		"name":"campaign",
		"namespace":"campaign.campaign_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "merchant_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
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
				"name":"budget",
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

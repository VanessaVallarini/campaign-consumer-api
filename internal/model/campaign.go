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
	Id         string  `avro:"id"`
	MerchantId string  `avro:"merchant_id"`
	Status     string  `avro:"status"`
	Lat        float64 `avro:"lat"`
	Long       float64 `avro:"long"`
	Budget     float64 `avro:"budget"`
	CreatedBy  string  `avro:"created_by"`
	UpdatedBy  string  `avro:"updated_by"`
}

const (
	CampaignAvro = `{
		"type":"record",
		"name":"campaign",
		"namespace":"campaign.campaign_value",
		"fields":[
			 {
				"name":"id",
				"type":"string"
			 },
			 {
				"name":"merchant_id",
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
				"name":"budget",
				"type":"double"
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

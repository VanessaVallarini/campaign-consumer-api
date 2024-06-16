package model

import (
	"github.com/google/uuid"
)

type ClickImpression struct {
	CampaignId uuid.UUID `json:"campaign_id"`
	SlugId     uuid.UUID `json:"slug_id"`
	UserId     uuid.UUID `json:"user_id"`
	EventType  string    `json:"event_type"`
	Lat        float64   `json:"lat"`
	Long       float64   `json:"long"`
}

type ClickImpressionEvent struct {
	CampaignId string  `avro:"campaign_id"`
	SlugId     string  `avro:"slug_id"`
	UserId     string  `avro:"user_id"`
	EventType  string  `avro:"event_type"`
	Lat        float64 `avro:"lat"`
	Long       float64 `avro:"long"`
}

const (
	ClickImpressionAvro = `{
		"type":"record",
		"name":"click_impression",
		"namespace":"campaign.campaign_click_impression_value",
		"fields":[
			 {
				"name":"campaign_id",
				"type":"string"
			 },
			 {
				"name":"slug_id",
				"type":"string"
			 },
			 {
				"name":"user_id",
				"type":"string"
			 },
			 {
				"name":"event_type",
				"type":"double"
			 },
			 {
				"name":"lat",
				"type":"double"
			 },
			 {
				"name":"long",
				"type":"double"
			 }   
		]
	 }`
)

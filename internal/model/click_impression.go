package model

import (
	"time"

	"github.com/google/uuid"
)

type ClickImpressionEvent struct {
	CampaignId uuid.UUID `avro:"campaign_id"`
	SlugId     uuid.UUID `avro:"slug_id"`
	UserId     uuid.UUID `avro:"user_id"`
	EventType  string    `avro:"event_type"`
	Lat        float64   `avro:"lat"`
	Long       float64   `avro:"long"`
	EventTime  time.Time `avro:"event_time"`
}

const (
	ClickImpressionAvro = `{
		"type":"record",
		"name":"click_impression",
		"namespace":"campaign.campaign_click_impression_value",
		"fields":[
			{
				"name": "campaign_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "slug_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "user_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name":"event_type",
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
				"name": "event_time",
				"type": {
				"type": "long",
				"logicalType": "timestamp-millis"
				}
			}
		]
	 }`
)

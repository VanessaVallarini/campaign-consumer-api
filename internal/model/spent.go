package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Spent struct {
	Id               uuid.UUID `json:"id"`
	CampaignId       uuid.UUID `json:"campaign_id"`
	MerchantId       uuid.UUID `json:"merchant_id"`
	Bucket           string    `json:"bucket"`
	TotalSpent       float64   `json:"total_spent"`
	TotalClicks      float64   `json:"total_clicks"`
	TotalImpressions float64   `json:"total_impressions"`
}

type EventType string

const (
	Click      EventType = "CLICK"
	Impression EventType = "IMPRESSION"
)

func ValidateEventType(eventType string) error {
	validEventType := map[string]bool{
		string(Click):      true,
		string(Impression): true,
	}

	if !validEventType[eventType] {

		return ErrInvalid
	}

	return nil
}

func EventTypeFromString(eventType string) EventType {
	switch eventType {
	case string(Click):
		return Click
	case string(Impression):
		return Impression
	default:
		return ""
	}
}

const (
	SpentEventAvro = `{
		"type":"record",
		"name":"spent",
		"namespace":"campaign.campaign_spent_value",
		"fields":[
			{
				"name": "campaign_id",
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
				"name": "session_id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
			},
			{
				"name": "slug_name",
				"type": "string"
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
				"name":"ip",
				"type":"string"
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

type SpentEvent struct {
	CampaignId uuid.UUID `json:"campaign_id" avro:"campaign_id"`
	MerchantId uuid.UUID `json:"merchant_id" avro:"merchant_id"`
	SessionId  uuid.UUID `json:"session_id" avro:"session_id"`
	SlugName   string    `json:"slug_name" avro:"slug_name"`
	UserId     uuid.UUID `json:"user_id" avro:"user_id"`
	EventType  string    `json:"event_type" avro:"event_type"`
	IP         string    `json:"ip" avro:"ip"`
	EventTime  time.Time `json:"event_time" avro:"event_time"`
}

func (se SpentEvent) ValidateSpentEvent() error {
	if se.CampaignId == (uuid.UUID{}) {

		return fmt.Errorf("invalid campaign id %v", se.CampaignId)
	}

	if se.MerchantId == (uuid.UUID{}) {

		return fmt.Errorf("invalid merchant id %v", se.MerchantId)
	}

	if se.SlugName == "" {

		return fmt.Errorf("invalid slug name %v", se.SlugName)
	}

	if se.UserId == (uuid.UUID{}) {

		return fmt.Errorf("invalid user id %v", se.UserId)
	}

	err := ValidateEventType(se.EventType)
	if err != nil {

		return fmt.Errorf("invalid event type %s", se.EventType)
	}

	if se.IP == "" {

		return fmt.Errorf("invalid ip %s", se.IP)
	}

	return nil
}

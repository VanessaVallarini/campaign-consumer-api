package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

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
				"name":"created_by",
				"type":"string"
			},
			{
				"name":"updated_by",
				"type":"string"
			},
			{
				"name": "created_at",
				"type": {
				"type": "long",
				"logicalType": "timestamp-millis"
				}
			},
			{
				"name": "updated_at",
				"type": {
				"type": "long",
				"logicalType": "timestamp-millis"
				}
			}		   
		]
	 }`
)

type Slug struct {
	Id        uuid.UUID `json:"id" avro:"id"`
	Name      string    `json:"name" avro:"name"`
	Status    string    `json:"status" avro:"status"`
	Cost      float64   `json:"cost" avro:"cost"`
	CreatedBy string    `json:"created_by" avro:"created_by"`
	UpdatedBy string    `json:"updated_by" avro:"updated_by"`
	CreatedAt time.Time `json:"created_at" avro:"created_at"`
	UpdatedAt time.Time `json:"updated_at" avro:"updated_at"`
}

func (s Slug) ValidateSlug() error {
	s.Name = strings.ToUpper(s.Name)
	s.Status = strings.ToUpper(s.Status)
	s.CreatedBy = strings.ToLower(s.CreatedBy)
	s.UpdatedBy = strings.ToLower(s.UpdatedBy)

	if s.Name == "" {

		return fmt.Errorf("invalid slug name %s", s.Name)
	}

	err := ValidateStatus(s.Status)
	if err != nil {

		return fmt.Errorf("invalid slug status %s", s.Status)
	}

	if s.Cost < 0 {

		return fmt.Errorf("invalid slug cost %f", s.Cost)
	}

	if s.CreatedBy == "" {

		return fmt.Errorf("invalid slug createdBy %s", s.CreatedBy)
	}

	if s.UpdatedBy == "" {

		return fmt.Errorf("invalid slug updatedBy %s", s.UpdatedBy)
	}

	return nil
}

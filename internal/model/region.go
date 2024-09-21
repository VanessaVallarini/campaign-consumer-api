package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

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

type Region struct {
	Id        uuid.UUID `json:"id" avro:"id"`
	Name      string    `json:"name" avro:"name"`
	Status    string    `json:"status" avro:"status"`
	Lat       float64   `json:"lat" avro:"lat"`
	Long      float64   `json:"long" avro:"long"`
	Cost      float64   `json:"cost" avro:"cost"`
	CreatedBy string    `json:"created_by" avro:"created_by"`
	UpdatedBy string    `json:"updated_by" avro:"updated_by"`
	CreatedAt time.Time `json:"created_at" avro:"created_at"`
	UpdatedAt time.Time `json:"updated_at" avro:"updated_at"`
}

func (r Region) ValidateRegion() error {
	r.Name = strings.ToUpper(r.Name)
	r.Status = strings.ToUpper(r.Status)
	r.CreatedBy = strings.ToLower(r.CreatedBy)
	r.UpdatedBy = strings.ToLower(r.UpdatedBy)

	if r.Name == "" {

		return fmt.Errorf("invalid region name %s", r.Name)
	}

	err := ValidateStatus(r.Status)
	if err != nil {

		return fmt.Errorf("invalid region status %s", r.Status)
	}

	if r.Lat < -90 || r.Lat > 90 {

		return fmt.Errorf("invalid region lat %f", r.Lat)
	}

	if r.Long < -90 || r.Long > 90 {

		return fmt.Errorf("invalid region long %f", r.Long)
	}

	if r.Cost < 0 {

		return fmt.Errorf("invalid region cost %f", r.Cost)
	}

	if r.CreatedBy == "" {

		return fmt.Errorf("invalid region createdBy %s", r.CreatedBy)
	}

	if r.UpdatedBy == "" {

		return fmt.Errorf("invalid region updatedBy %s", r.UpdatedBy)
	}

	return nil
}

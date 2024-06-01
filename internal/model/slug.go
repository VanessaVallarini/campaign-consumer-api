package model

import (
	"time"

	"github.com/google/uuid"
)

type Slug struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Cost      float64   `json:"cost"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SlugEvent struct {
	Id        string  `json:"id" avro:"id"`
	Name      string  `json:"name" avro:"name"`
	Active    string  `json:"active" avro:"active"` //conseguimos não usar apenas string?
	Cost      float64 `json:"cost" avro:"cost"`
	CreatedBy string  `json:"created_by" avro:"created_by"`
	UpdatedBy string  `json:"updated_by" avro:"updated_by"`
}

const (
	SlugAvro = `{
		"type":"record",
		"name":"slug",
		"namespace":"campaign.campaign_slug_value",
		"fields":[
			 {
				"name":"id",
				"type":"string"
			 },
			 {
				"name":"name",
				"type":"string"
			 },
			 {
				"name":"active",
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
			 }		   
		]
	 }`
)

package model

import (
	"time"

	"github.com/google/uuid"
)

type Owner struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OwnerEvent struct {
	Id        string `json:"id" avro:"id"`
	Email     string `json:"email" avro:"email"`
	Active    string `json:"active" avro:"active"` //conseguimos não usar apenas string?
	CreatedBy string `json:"created_by" avro:"created_by"`
	UpdatedBy string `json:"updated_by" avro:"updated_by"`
}

const (
	OwnerAvro = `{
		"type":"record",
		"name":"owner",
		"namespace":"campaign.campaign_owner_value",
		"fields":[
			 {
				"name":"id",
				"type":"string"
			 },
			 {
				"name":"email",
				"type":"string"
			 },
			 {
				"name":"active",
				"type":"string"
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

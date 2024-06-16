package model

import (
	"time"

	"github.com/google/uuid"
)

type OwnerStatus string

const (
	ActiveOwner   OwnerStatus = "ACTIVE"
	InactiveOwner OwnerStatus = "INACTIVE"
)

type Owner struct {
	Id        uuid.UUID   `json:"id"`
	Email     string      `json:"email"`
	Status    OwnerStatus `json:"status"`
	CreatedBy string      `json:"created_by"`
	UpdatedBy string      `json:"updated_by"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OwnerEvent struct {
	Id        string `json:"id" avro:"id"`
	Email     string `json:"email" avro:"email"`
	Status    string `json:"status" avro:"status"`
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
				"name":"status",
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

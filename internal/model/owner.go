package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	OwnerAvro = `{
		"type":"record",
		"name":"owner",
		"namespace":"campaign.campaign_owner_value",
		"fields":[
			{
				"name": "id",
				"type": {
				"type": "string",
				"logicalType": "UUID"
				}
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

type Owner struct {
	Id        uuid.UUID `json:"id" avro:"id"`
	Email     string    `json:"email" avro:"email"`
	Status    string    `json:"status" avro:"status"`
	CreatedBy string    `json:"created_by" avro:"created_by"`
	UpdatedBy string    `json:"updated_by" avro:"updated_by"`
	CreatedAt time.Time `json:"created_at" avro:"created_at"`
	UpdatedAt time.Time `json:"updated_at" avro:"updated_at"`
}

func (o Owner) ValidateOwner() error {
	o.Email = strings.ToUpper(o.Email)
	o.Status = strings.ToUpper(o.Status)
	o.CreatedBy = strings.ToLower(o.CreatedBy)
	o.UpdatedBy = strings.ToLower(o.UpdatedBy)

	err := validateEmail(o.Email)
	if err != nil {

		return err
	}

	err = ValidateStatus(o.Status)
	if err != nil {

		return fmt.Errorf("invalid owner status %s", o.Status)
	}

	if o.CreatedBy == "" {

		return fmt.Errorf("invalid owner createdBy %s", o.CreatedBy)
	}

	if o.UpdatedBy == "" {

		return fmt.Errorf("invalid owner updatedBy %s", o.UpdatedBy)
	}

	return nil
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + `-{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

func validateEmail(email string) error {
	// Check for empty email
	if email == "" {
		return fmt.Errorf("invalid owner email: empty")
	}

	// Optional: More stringent regex validation
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid owner email format: stricter validation failed")
	}

	return nil
}

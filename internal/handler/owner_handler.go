package handler

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka"
	easyzap "github.com/lockp111/go-easyzap"
)

type OwnerProcessor interface {
	OwnerProcessor(model.OwnerEvent) error
}

func MakeOwnerEventHandler(processor OwnerProcessor) func(msg *sarama.ConsumerMessage, srClient kafka.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient kafka.SchemaRegistryClient, subject string) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		// Decode msg.Value into model.Owner
		var owner model.OwnerEvent
		if err := srClient.Decode(msg.Value, &owner, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update owner")

			return err
		}

		easyzap.Infof("got owner event for %s", owner.Email)
		if err := processor.OwnerProcessor(owner); err != nil {

			return err
		}

		return nil
	}
}

package handler

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	easyzap "github.com/lockp111/go-easyzap"
)

type OwnerService interface {
	Upsert(context.Context, model.Owner) error
}

func MakeOwnerEventHandler(ownerService OwnerService) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		// Decode msg.Value into model.Owner
		var owner model.Owner
		if err := srClient.Decode(msg.Value, &owner, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update owner")

			return err
		}

		easyzap.Infof("got owner event for %s", owner.Email)
		if err := ownerService.Upsert(context.Background(), owner); err != nil {

			return err
		}

		return nil
	}
}

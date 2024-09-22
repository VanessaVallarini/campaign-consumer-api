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
		ctx := context.Background()

		if msg == nil {
			err := errors.New("Invalid message pointer")
			easyzap.Error(ctx, err)

			return err
		}

		var owner model.Owner
		if err := srClient.Decode(msg.Value, &owner, subject); err != nil {
			easyzap.Error(ctx, err, "error during decode message consumer kafka on create or update owner")

			return err
		}

		if err := ownerService.Upsert(ctx, owner); err != nil {

			return err
		}

		return nil
	}
}

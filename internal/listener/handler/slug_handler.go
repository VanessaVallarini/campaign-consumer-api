package handler

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	easyzap "github.com/lockp111/go-easyzap"
)

type SlugService interface {
	Upsert(context.Context, model.Slug) error
}

func MakeSlugEventHandler(slugService SlugService) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		ctx := context.Background()

		if msg == nil {
			err := errors.New("Invalid message pointer")
			easyzap.Error(ctx, err)

			return err
		}

		var slug model.Slug
		if err := srClient.Decode(msg.Value, &slug, subject); err != nil {
			easyzap.Error(ctx, err, "error during decode message consumer kafka on create or update slug")

			return err
		}

		if err := slugService.Upsert(ctx, slug); err != nil {

			return err
		}

		return nil
	}
}

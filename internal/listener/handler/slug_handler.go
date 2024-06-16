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
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		// Decode msg.Value into model.SlugEvent
		var slug model.Slug
		if err := srClient.Decode(msg.Value, &slug, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update slug")

			return err
		}

		easyzap.Infof("got slug event for %s", slug.Name)
		if err := slugService.Upsert(context.Background(), slug); err != nil {

			return err
		}

		return nil
	}
}

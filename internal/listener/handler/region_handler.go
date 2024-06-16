package handler

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	easyzap "github.com/lockp111/go-easyzap"
)

type RegionService interface {
	Upsert(context.Context, model.Region) error
}

func MakeRegionEventHandler(regionService RegionService) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		// Decode msg.Value into model.Region
		var region model.Region
		if err := srClient.Decode(msg.Value, &region, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update region")

			return err
		}

		easyzap.Infof("got region event for %s", region.Name)
		if err := regionService.Upsert(context.Background(), region); err != nil {

			return err
		}

		return nil
	}
}

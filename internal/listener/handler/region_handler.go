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
		ctx := context.Background()

		if msg == nil {
			err := errors.New("Invalid message pointer")
			easyzap.Error(ctx, err)

			return err
		}

		var region model.Region
		if err := srClient.Decode(msg.Value, &region, subject); err != nil {
			easyzap.Error(ctx, err, "error during decode message consumer kafka on create or update region")

			return err
		}

		if err := regionService.Upsert(ctx, region); err != nil {

			return err
		}

		return nil
	}
}

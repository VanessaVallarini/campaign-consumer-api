package handler

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	easyzap "github.com/lockp111/go-easyzap"
)

type MerchantService interface {
	Upsert(context.Context, model.Merchant) error
}

func MakeMerchantEventHandler(merchantService MerchantService) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		ctx := context.Background()

		if msg == nil {
			err := errors.New("Invalid message pointer")
			easyzap.Error(ctx, err)

			return err
		}

		var merchant model.Merchant
		if err := srClient.Decode(msg.Value, &merchant, subject); err != nil {
			easyzap.Error(ctx, err, "error during decode message consumer kafka on create or update merchant")

			return err
		}

		if err := merchantService.Upsert(ctx, merchant); err != nil {

			return err
		}

		return nil
	}
}

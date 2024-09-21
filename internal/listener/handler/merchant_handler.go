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
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		var merchant model.Merchant
		if err := srClient.Decode(msg.Value, &merchant, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update merchant")

			return err
		}

		if err := merchantService.Upsert(context.Background(), merchant); err != nil {

			return err
		}

		return nil
	}
}

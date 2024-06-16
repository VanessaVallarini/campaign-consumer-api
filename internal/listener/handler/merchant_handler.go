package handler

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	easyzap "github.com/lockp111/go-easyzap"
)

type MerchantProcessor interface {
	MerchantProcessor(model.MerchantEvent) error
}

func MakeMerchantEventHandler(processor MerchantProcessor) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		// Decode msg.Value into model.MerchantEvent
		var merchantEvent model.MerchantEvent
		if err := srClient.Decode(msg.Value, &merchantEvent, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update merchant")

			return err
		}

		easyzap.Infof("got merchant event for %s", merchantEvent.Name)
		if err := processor.MerchantProcessor(merchantEvent); err != nil {

			return err
		}

		return nil
	}
}

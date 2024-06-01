package handler

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	easyzap "github.com/lockp111/go-easyzap"
)

type SlugProcessor interface {
	SlugProcessor(model.SlugEvent) error
}

func MakeSlugEventHandler(processor SlugProcessor) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		// Decode msg.Value into model.SlugEvent
		var slugEvent model.SlugEvent
		if err := srClient.Decode(msg.Value, &slugEvent, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update slug")

			return err
		}

		easyzap.Infof("got slug event for %s", slugEvent.Name)
		if err := processor.SlugProcessor(slugEvent); err != nil {

			return err
		}

		return nil
	}
}

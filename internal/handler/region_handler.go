package handler

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	easyzap "github.com/lockp111/go-easyzap"
)

type RegionProcessor interface {
	RegionProcessor(model.RegionEvent) error
}

func MakeRegionEventHandler(processor RegionProcessor) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		// Decode msg.Value into model.RegionEvent
		var regionEvent model.RegionEvent
		if err := srClient.Decode(msg.Value, &regionEvent, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update region")

			return err
		}

		easyzap.Infof("got region event for %s", regionEvent.Name)
		if err := processor.RegionProcessor(regionEvent); err != nil {

			return err
		}

		return nil
	}
}

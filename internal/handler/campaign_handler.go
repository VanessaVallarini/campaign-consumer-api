package handler

import (
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/client"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	easyzap "github.com/lockp111/go-easyzap"
)

type CampaignProcessor interface {
	CampaignProcessor(model.CampaignEvent) error
}

func MakeCampaignEventHandler(processor CampaignProcessor) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		// Decode msg.Value into model.CampaignEvent
		var campaignEvent model.CampaignEvent
		if err := srClient.Decode(msg.Value, &campaignEvent, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update campaign")

			return err
		}

		easyzap.Infof("got campaign event for merchant id %s", campaignEvent.MerchantId)
		if err := processor.CampaignProcessor(campaignEvent); err != nil {

			return err
		}

		return nil
	}
}

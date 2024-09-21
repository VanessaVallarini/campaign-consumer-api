package handler

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	easyzap "github.com/lockp111/go-easyzap"
)

type CampaignService interface {
	Upsert(context.Context, model.Campaign) error
}

func MakeCampaignEventHandler(campaignService CampaignService) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("Invalid message pointer")
		}

		var campaign model.Campaign
		if err := srClient.Decode(msg.Value, &campaign, subject); err != nil {
			easyzap.Error(err, "error during decode message consumer kafka on create or update campaign")

			return err
		}

		if err := campaignService.Upsert(context.Background(), campaign); err != nil {

			return err
		}

		return nil
	}
}

package handler

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
	easyzap "github.com/lockp111/go-easyzap"
)

type SpentProcessor interface {
	ProcessSpentEvent(context.Context, model.SpentEvent) error
}

func MakeSpentEventHandler(spentProcessor SpentProcessor) func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
	return func(msg *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error {
		ctx := context.Background()

		if msg == nil {
			err := errors.New("Invalid message pointer")
			easyzap.Error(ctx, err)

			return err
		}

		var spent model.SpentEvent
		if err := srClient.Decode(msg.Value, &spent, subject); err != nil {
			easyzap.Error(ctx, err, "error during decode message consumer kafka on create or update spent")

			return err
		}

		if err := spentProcessor.ProcessSpentEvent(ctx, spent); err != nil {

			return err
		}

		return nil
	}
}

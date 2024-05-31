package handler

import (
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	easyzap "github.com/lockp111/go-easyzap"
)

type OwnerProcessor interface {
	OwnerProcessor(model.OwnerEvent) error
}

func MakeOwnerEventHandler(processor OwnerProcessor) func(msg *sarama.ConsumerMessage, obj interface{}) error {
	return func(msg *sarama.ConsumerMessage, obj interface{}) error {
		if msg == nil {
			easyzap.Error("invalid message pointer")

			return errors.New("invalid message pointer")
		}

		// Decode msg.Value into model.Owner
		var message model.OwnerEvent
		err := json.Unmarshal(msg.Value, &message)
		if err != nil {
			easyzap.Errorf("Error decoding JSON: %v", err)

			return err
		}

		easyzap.Infof("got owner event for %s", message.Email)
		err = processor.OwnerProcessor(message)
		if err != nil {
			return err
		}

		return nil
	}
}

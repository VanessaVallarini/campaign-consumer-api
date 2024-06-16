package consumer

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	easyzap "github.com/lockp111/go-easyzap"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/kafka/client"
)

type MessageHandler func(message *sarama.ConsumerMessage, srClient client.SchemaRegistryClient, subject string) error

type Consumer struct {
	ready               chan bool
	kafkaClient         sarama.Client
	saramaClusterAdmin  sarama.ClusterAdmin
	saramaConsumerGroup sarama.ConsumerGroup
	saramaConfig        *sarama.Config
	srClient            client.SchemaRegistryClient
	messageHandler      MessageHandler
	subject             string
}

func NewConsumer(ctx context.Context, brokerConfig config.KafkaConfig, srClient client.SchemaRegistryClient, messageHandler MessageHandler) Consumer {
	kafkaClient, saramaClusterAdmin, saramaConsumerGroup, saramaConfig := client.NewKafkaClient(brokerConfig)

	return Consumer{
		ready:               make(chan bool),
		kafkaClient:         kafkaClient,
		saramaClusterAdmin:  saramaClusterAdmin,
		saramaConsumerGroup: saramaConsumerGroup,
		saramaConfig:        saramaConfig,
		srClient:            srClient,
		messageHandler:      messageHandler,
		subject:             brokerConfig.Subject,
	}
}

func (c Consumer) ConsumerStart(brokerConfig config.KafkaConfig) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ctx := context.Background()
			handler := sarama.ConsumerGroupHandler(&c)
			if err := c.saramaConsumerGroup.Consume(ctx, []string{brokerConfig.Topic}, handler); err != nil {
				easyzap.Error(err, "failed create consumer")
			}
			if ctx.Err() != nil {
				err := c.saramaConsumerGroup.Close()
				if err != nil {
					easyzap.Error(err, "failed to close consumer")
				}

				easyzap.Info("consumer closed, consuming again")
			}

			c.ready = make(chan bool)
		}
	}()

	<-c.ready
	easyzap.Info("consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigterm:
		easyzap.Info("terminating: via signal")
	}

	wg.Wait()
	if err := c.saramaConsumerGroup.Close(); err != nil {
		easyzap.Fatalf("Error closing groupClient: %v", err)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			if err := c.messageHandler(message, c.srClient, c.subject); err != nil {
				return err
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

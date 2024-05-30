package client

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	easyzap "github.com/lockp111/go-easyzap"
)

const _saramaTimeoutFlushMs = 500

type KafkaClient struct {
	saramaClient       sarama.Client
	saramaClusterAdmin sarama.ClusterAdmin
}

func NewKafkaClient(ctx context.Context, cfg config.KafkaConfig) KafkaClient {
	saramaConfig := generateSaramaConfig(cfg)

	saramaClient, err := sarama.NewClient(cfg.Brokers, saramaConfig)
	if err != nil {
		easyzap.Fatal(ctx, err, "kafka client failed to new kafka client")
	}

	saramaClusterAdmin, err := sarama.NewClusterAdminFromClient(saramaClient)
	if err != nil {
		easyzap.Fatal(ctx, err, "kafka client failed to new cluster admin")
	}

	if !cfg.UseAuthentication {
		createTopic(cfg.Topic, saramaClusterAdmin)
	}

	return KafkaClient{
		saramaClient:       saramaClient,
		saramaClusterAdmin: saramaClusterAdmin,
	}
}

func generateSaramaConfig(cfg config.KafkaConfig) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaTimeout := cfg.Timeout * time.Millisecond

	saramaConfig.ClientID = cfg.ClientId
	saramaConfig.Version = sarama.V3_0_0_0
	saramaConfig.Net.DialTimeout = saramaTimeout
	saramaConfig.Net.ReadTimeout = saramaTimeout
	saramaConfig.Net.WriteTimeout = saramaTimeout
	saramaConfig.Metadata.Timeout = saramaTimeout
	saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	saramaConfig.Producer.Flush.Frequency = _saramaTimeoutFlushMs * time.Millisecond
	saramaConfig.Metadata.Retry.Max = cfg.RetryMax
	saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{cfg.BalanceStrategy}

	if cfg.UseAuthentication {
		saramaConfig.Net.SASL.Mechanism = sarama.SASLMechanism(cfg.SaslMechanism)
		saramaConfig.Net.SASL.User = cfg.User
		saramaConfig.Net.SASL.Password = cfg.Password
		saramaConfig.Net.TLS.Enable = cfg.EnableTLS
		saramaConfig.Net.SASL.Enable = true
		setAuthentication(saramaConfig)
	}

	return saramaConfig
}

func setAuthentication(conf *sarama.Config) {
	switch conf.Net.SASL.Mechanism {
	case sarama.SASLTypeSCRAMSHA512:
		scram512Fn := func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		conf.Net.SASL.SCRAMClientGeneratorFunc = scram512Fn
	case sarama.SASLTypeSCRAMSHA256:
		scram256Fn := func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		conf.Net.SASL.SCRAMClientGeneratorFunc = scram256Fn
	}
}

func createTopic(topic string, saramaClusterAdmin sarama.ClusterAdmin) error {
	err := saramaClusterAdmin.CreateTopic(topic,
		&sarama.TopicDetail{
			NumPartitions:     4,
			ReplicationFactor: 1,
		},
		false)
	if err != nil {
		easyzap.Warnf("kafka client failed create topic: %s. details: %v", topic, err)
	}
	return nil
}

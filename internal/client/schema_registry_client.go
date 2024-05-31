package client

import (
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"github.com/riferrei/srclient"
)

func NewSchemaRegistry(brokerConfig config.KafkaConfig) *srclient.SchemaRegistryClient {
	srClient := srclient.CreateSchemaRegistryClient(brokerConfig.SchemaRegistryConfig.Host)
	if brokerConfig.UseAuthentication {
		srClient.SetCredentials(brokerConfig.SchemaRegistryConfig.User, brokerConfig.SchemaRegistryConfig.Password)
	}

	return srClient
}

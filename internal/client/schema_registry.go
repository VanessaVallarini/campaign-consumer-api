package client

import (
	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"github.com/riferrei/srclient"
)

type SchemaRegistry struct {
	*srclient.SchemaRegistryClient
}

func NewSchemaRegistry(cfg config.KafkaConfig) SchemaRegistry {
	src := srclient.CreateSchemaRegistryClient(cfg.SchemaRegistryConfig.Host)
	if cfg.UseAuthentication {
		src.SetCredentials(cfg.SchemaRegistryConfig.User, cfg.SchemaRegistryConfig.Password)
	}

	return SchemaRegistry{src}
}

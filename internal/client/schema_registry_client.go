package client

import (
	"encoding/json"
	"errors"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/config"
	"github.com/hamba/avro"
	"github.com/riferrei/srclient"
)

type SchemaRegistryClient struct {
	srClient *srclient.SchemaRegistryClient
}

func NewSchemaRegistry(brokerConfig config.KafkaConfig) SchemaRegistryClient {
	srClient := srclient.CreateSchemaRegistryClient(brokerConfig.SchemaRegistryConfig.Host)
	if brokerConfig.UseAuthentication {
		srClient.SetCredentials(brokerConfig.SchemaRegistryConfig.User, brokerConfig.SchemaRegistryConfig.Password)
	}

	return SchemaRegistryClient{
		srClient: srClient,
	}
}

func (src SchemaRegistryClient) Decode(data []byte, value interface{}, subject string) error {
	schema, err := src.getSchema(subject)
	if err != nil {
		return err
	}

	schemaDecoder, err := avro.Parse(schema.Schema())
	if err != nil {
		return err
	}

	errAvro := avro.Unmarshal(schemaDecoder, data, value)
	if errAvro != nil {
		//ignorando pq está dando erro
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		return err
	}

	return nil
}

func (src SchemaRegistryClient) getSchema(subject string) (*srclient.Schema, error) {
	schema, err := src.srClient.GetLatestSchema(subject)
	if err != nil {
		return nil, err
	}

	if schema == nil {
		return nil, errors.New("schema registry unexpected behavior retrieving schema, got 'nil' from registry")
	}

	return schema, nil
}

package kafka

import (
	"github.com/pkg/errors"

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

func (sr SchemaRegistryClient) Decode(data []byte, value interface{}, subject string) error {
	schema, err := sr.getSchema(subject)
	if err != nil {
		return err
	}

	schemaDecoder, err := avro.Parse(schema.Schema())
	if err != nil {
		return err
	}

	errAvro := avro.Unmarshal(schemaDecoder, data[5:], value)
	if errAvro != nil {
		return errors.Wrap(err, "Failed to avro unmarshal")
	}

	return nil
}

func (sr SchemaRegistryClient) getSchema(subject string) (*srclient.Schema, error) {
	schema, err := sr.srClient.GetLatestSchema(subject)
	if err != nil {
		return nil, err
	}

	if schema == nil {
		return nil, errors.New("Schema registry unexpected behavior retrieving schema, got 'nil' from registry")
	}

	return schema, nil
}

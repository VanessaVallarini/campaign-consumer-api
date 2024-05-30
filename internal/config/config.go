package config

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	easyzap "github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	AppName              string
	ServerHost           string
	MetaHost             string
	TimeLocation         string
	Database             DatabaseConfig
	KafkaOwner           KafkaConfig
	KafkaSlug            KafkaConfig
	KafkaMerchant        KafkaConfig
	KafkaCampaign        KafkaConfig
	KafkaClickImpression KafkaConfig
}

type DatabaseConfig struct {
	Host     string
	Username string
	Password string
	Database string
	Port     int
	Conn     Conn
}

type Conn struct {
	Min      int
	Max      int
	Lifetime string
	IdleTime string
}

type KafkaConfig struct {
	ClientId          string
	ConsumerGroupId   string
	Brokers           []string
	Acks              string
	Timeout           time.Duration
	UseAuthentication bool
	EnableTLS         bool
	EnableEvents      bool
	SaslMechanism     string
	User              string
	Password          string
	RetryMax          int
	Topic             string
	Subject           string
	BalanceStrategy   sarama.BalanceStrategy
	SchemaRegistryConfig
}

type SchemaRegistryConfig struct {
	Host     string
	User     string
	Password string
}

var (
	onceConfigs sync.Once
	configInit  sync.Once
	config      Config
	viperCfg    = viper.New()
)

func initConfig() {
	viperCfg.AddConfigPath("internal/config/")
	viperCfg.SetConfigName("configuration")
	viperCfg.SetConfigType("yml")

	setConfigDefaults()

	if err := viperCfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			err := errors.Wrapf(err, "Error reading config file: %s", err)
			easyzap.Fatal(context.Background(), err, "Unable to keep the service without config file")
		}
	}

	viperCfg.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viperCfg.AutomaticEnv()
}

func setConfigDefaults() {
	viperCfg.SetDefault("app.name", "campaign-consumer-api")
	viperCfg.SetDefault("server.host", "0.0.0.0:8080")
	viperCfg.SetDefault("meta.host", "0.0.0.0:8081")
}

func GetConfig() Config {
	configInit.Do(initConfig)
	onceConfigs.Do(func() {
		config = Config{
			AppName:      viperCfg.GetString("app.name"),
			ServerHost:   viperCfg.GetString("server.host"),
			MetaHost:     viperCfg.GetString("meta.host"),
			TimeLocation: viperCfg.GetString("time-location"),
			Database: DatabaseConfig{
				Host:     viperCfg.GetString("database.host"),
				Username: viperCfg.GetString("database.username"),
				Password: viperCfg.GetString("database.password"),
				Database: viperCfg.GetString("database.database"),
				Port:     viperCfg.GetInt("database.port"),
				Conn: Conn{
					Min:      viperCfg.GetInt("database.conn.min"),
					Max:      viperCfg.GetInt("database.conn.max"),
					Lifetime: viperCfg.GetString("database.conn.lifetime"),
					IdleTime: viperCfg.GetString("database.conn.idletime"),
				},
			},
			KafkaOwner: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka-owner.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka-owner.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka-owner.brokers"),
				Acks:              viperCfg.GetString("kafka-owner.acks"),
				Timeout:           viperCfg.GetDuration("kafka-owner.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka-owner.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka-owner.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka-owner.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka-owner.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-owner.user"),
				Password:          viperCfg.GetString("kafka-owner.password"),
				RetryMax:          viperCfg.GetInt("kafka-owner.retry-max"),
				Topic:             viperCfg.GetString("kafka-owner.topic"),
				Subject:           viperCfg.GetString("kafka-owner.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka-owner.schema-registry.host"),
					User:     viperCfg.GetString("kafka-owner.schema-registry.user"),
					Password: viperCfg.GetString("kafka-owner.schema-registry.password"),
				},
			},
			KafkaSlug: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka-slug.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka-slug.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka-slug.brokers"),
				Acks:              viperCfg.GetString("kafka-slug.acks"),
				Timeout:           viperCfg.GetDuration("kafka-slug.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka-slug.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka-slug.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka-slug.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka-slug.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-slug.user"),
				Password:          viperCfg.GetString("kafka-slug.password"),
				RetryMax:          viperCfg.GetInt("kafka-slug.retry-max"),
				Topic:             viperCfg.GetString("kafka-slug.topic"),
				Subject:           viperCfg.GetString("kafka-slug.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka-slug.schema-registry.host"),
					User:     viperCfg.GetString("kafka-slug.schema-registry.user"),
					Password: viperCfg.GetString("kafka-slug.schema-registry.password"),
				},
			},
			KafkaMerchant: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka-merchant.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka-merchant.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka-merchant.brokers"),
				Acks:              viperCfg.GetString("kafka-merchant.acks"),
				Timeout:           viperCfg.GetDuration("kafka-merchant.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka-merchant.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka-merchant.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka-merchant.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka-merchant.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-merchant.user"),
				Password:          viperCfg.GetString("kafka-merchant.password"),
				RetryMax:          viperCfg.GetInt("kafka-merchant.retry-max"),
				Topic:             viperCfg.GetString("kafka-merchant.topic"),
				Subject:           viperCfg.GetString("kafka-merchant.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka-merchant.schema-registry.host"),
					User:     viperCfg.GetString("kafka-merchant.schema-registry.user"),
					Password: viperCfg.GetString("kafka-merchant.schema-registry.password"),
				},
			},
			KafkaCampaign: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka-campaign.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka-campaign.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka-campaign.brokers"),
				Acks:              viperCfg.GetString("kafka-campaign.acks"),
				Timeout:           viperCfg.GetDuration("kafka-campaign.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka-campaign.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka-campaign.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka-campaign.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka-campaign.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-campaign.user"),
				Password:          viperCfg.GetString("kafka-campaign.password"),
				RetryMax:          viperCfg.GetInt("kafka-campaign.retry-max"),
				Topic:             viperCfg.GetString("kafka-campaign.topic"),
				Subject:           viperCfg.GetString("kafka-campaign.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka-campaign.schema-registry.host"),
					User:     viperCfg.GetString("kafka-campaign.schema-registry.user"),
					Password: viperCfg.GetString("kafka-campaign.schema-registry.password"),
				},
			},
			KafkaClickImpression: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka-click-impression.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka-click-impression.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka-click-impression.brokers"),
				Acks:              viperCfg.GetString("kafka-click-impression.acks"),
				Timeout:           viperCfg.GetDuration("kafka-click-impression.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka-click-impression.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka-click-impression.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka-click-impression.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka-click-impression.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-click-impression.user"),
				Password:          viperCfg.GetString("kafka-click-impression.password"),
				RetryMax:          viperCfg.GetInt("kafka-click-impression.retry-max"),
				Topic:             viperCfg.GetString("kafka-click-impression.topic"),
				Subject:           viperCfg.GetString("kafka-click-impression.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka-click-impression.schema-registry.host"),
					User:     viperCfg.GetString("kafka-click-impression.schema-registry.user"),
					Password: viperCfg.GetString("kafka-click-impression.schema-registry.password"),
				},
			},
		}
	})

	return config
}

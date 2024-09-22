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
	KafkaRegion          KafkaConfig
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
			easyzap.Fatal(context.Background(), err, "unable to keep the service without config file")
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
				ClientId:          viperCfg.GetString("kafka.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka.brokers"),
				Acks:              viperCfg.GetString("kafka.acks"),
				Timeout:           viperCfg.GetDuration("kafka.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-owner.user"),
				Password:          viperCfg.GetString("kafka-owner.password"),
				RetryMax:          viperCfg.GetInt("kafka.retry-max"),
				Topic:             viperCfg.GetString("kafka-owner.topic"),
				Subject:           viperCfg.GetString("kafka-owner.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka.schema-registry.host"),
					User:     viperCfg.GetString("kafka-owner.schema-registry.user"),
					Password: viperCfg.GetString("kafka-owner.schema-registry.password"),
				},
			},
			KafkaSlug: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka.brokers"),
				Acks:              viperCfg.GetString("kafka.acks"),
				Timeout:           viperCfg.GetDuration("kafka.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-slug.user"),
				Password:          viperCfg.GetString("kafka-slug.password"),
				RetryMax:          viperCfg.GetInt("kafka.retry-max"),
				Topic:             viperCfg.GetString("kafka-slug.topic"),
				Subject:           viperCfg.GetString("kafka-slug.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka.schema-registry.host"),
					User:     viperCfg.GetString("kafka-slug.schema-registry.user"),
					Password: viperCfg.GetString("kafka-slug.schema-registry.password"),
				},
			},
			KafkaRegion: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka.brokers"),
				Acks:              viperCfg.GetString("kafka.acks"),
				Timeout:           viperCfg.GetDuration("kafka.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-region.user"),
				Password:          viperCfg.GetString("kafka-region.password"),
				RetryMax:          viperCfg.GetInt("kafka.retry-max"),
				Topic:             viperCfg.GetString("kafka-region.topic"),
				Subject:           viperCfg.GetString("kafka-region.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka.schema-registry.host"),
					User:     viperCfg.GetString("kafka-region.schema-registry.user"),
					Password: viperCfg.GetString("kafka-region.schema-registry.password"),
				},
			},
			KafkaMerchant: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka.brokers"),
				Acks:              viperCfg.GetString("kafka.acks"),
				Timeout:           viperCfg.GetDuration("kafka.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-merchant.user"),
				Password:          viperCfg.GetString("kafka-merchant.password"),
				RetryMax:          viperCfg.GetInt("kafka.retry-max"),
				Topic:             viperCfg.GetString("kafka-merchant.topic"),
				Subject:           viperCfg.GetString("kafka-merchant.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka.schema-registry.host"),
					User:     viperCfg.GetString("kafka-merchant.schema-registry.user"),
					Password: viperCfg.GetString("kafka-merchant.schema-registry.password"),
				},
			},
			KafkaCampaign: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka.brokers"),
				Acks:              viperCfg.GetString("kafka.acks"),
				Timeout:           viperCfg.GetDuration("kafka.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-campaign.user"),
				Password:          viperCfg.GetString("kafka-campaign.password"),
				RetryMax:          viperCfg.GetInt("kafka.retry-max"),
				Topic:             viperCfg.GetString("kafka-campaign.topic"),
				Subject:           viperCfg.GetString("kafka-campaign.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka.schema-registry.host"),
					User:     viperCfg.GetString("kafka-campaign.schema-registry.user"),
					Password: viperCfg.GetString("kafka-campaign.schema-registry.password"),
				},
			},
			KafkaClickImpression: KafkaConfig{
				ClientId:          viperCfg.GetString("kafka.client-id"),
				ConsumerGroupId:   viperCfg.GetString("kafka.consumer-group-id"),
				Brokers:           viperCfg.GetStringSlice("kafka.brokers"),
				Acks:              viperCfg.GetString("kafka.acks"),
				Timeout:           viperCfg.GetDuration("kafka.timeout"),
				UseAuthentication: viperCfg.GetBool("kafka.use-authentication"),
				EnableTLS:         viperCfg.GetBool("kafka.enable-tls"),
				EnableEvents:      viperCfg.GetBool("kafka.enable-events"),
				SaslMechanism:     viperCfg.GetString("kafka.sasl-mechanism"),
				User:              viperCfg.GetString("kafka-click-impression.user"),
				Password:          viperCfg.GetString("kafka-click-impression.password"),
				RetryMax:          viperCfg.GetInt("kafka.retry-max"),
				Topic:             viperCfg.GetString("kafka-click-impression.topic"),
				Subject:           viperCfg.GetString("kafka-click-impression.subject"),
				BalanceStrategy:   sarama.NewBalanceStrategyRoundRobin(),
				SchemaRegistryConfig: SchemaRegistryConfig{
					Host:     viperCfg.GetString("kafka.schema-registry.host"),
					User:     viperCfg.GetString("kafka-click-impression.schema-registry.user"),
					Password: viperCfg.GetString("kafka-click-impression.schema-registry.password"),
				},
			},
		}
	})

	return config
}

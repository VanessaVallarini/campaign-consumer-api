package config

import (
	"context"
	"strings"
	"sync"

	easyzap "github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	AppName      string
	ServerHost   string
	MetaHost     string
	TimeLocation string
	Database     DatabaseConfig
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
		}
	})

	return config
}

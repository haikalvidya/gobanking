package config

import (
	"flag"
	"fmt"
	"gobanking/pkg/logger"
	"gobanking/pkg/mysql"
	"gobanking/pkg/nats"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Search microservice config path")
}

type Config struct {
	ServiceName string           `mapstructure:"serviceName"`
	Logger      logger.LogConfig `mapstructure:"logger"`
	Timeouts    Timeouts         `mapstructure:"timeouts" validate:"required"`
	Http        Http             `mapstructure:"http"`
	Mysql       *mysql.Config    `mapstructure:"mysql"`
	Nats        *nats.NatsConfig `mapstructure:"nats"`
}

type Timeouts struct {
	MysqlInitMilliseconds int  `mapstructure:"mysqlInitMilliseconds" validate:"required"`
	MysqlInitRetryCount   uint `mapstructure:"mysqlInitRetryCount" validate:"required"`
}

type Http struct {
	Port                string `mapstructure:"port" validate:"required"`
	Development         bool   `mapstructure:"development"`
	HttpClientDebug     bool   `mapstructure:"httpClientDebug"`
	DebugErrorsResponse bool   `mapstructure:"debugErrorsResponse"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv("CONFIG_PATH")
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, nil
}

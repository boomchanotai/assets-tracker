package config

import (
	"github.com/boomchanotai/assets-tracker/server/pkg/logger"
	"github.com/boomchanotai/assets-tracker/server/pkg/postgres"
	"github.com/spf13/viper"
)

const (
	defaultPath = "./server/apps/api"
)

type AppConfig struct {
	Name     string          `mapstructure:"name"`
	Port     int             `mapstructure:"port"`
	Logger   logger.Config   `mapstructure:"logger"`
	Postgres postgres.Config `mapstructure:"postgres"`
}

func Load() *AppConfig {
	appConfig := &AppConfig{}
	viper.SetConfigName("config")
	viper.AddConfigPath(defaultPath)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&appConfig); err != nil {
		panic("error occurs while unmarshalling the config")
	}

	return appConfig
}

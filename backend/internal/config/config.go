package config

import (
	"github.com/boomchanotai/assets-tracker/pkg/logger"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Logger logger.Config `mapstructure:"logger"`
}

func Load() *AppConfig {
	appConfig := &AppConfig{}
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
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

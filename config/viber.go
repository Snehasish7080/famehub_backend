package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	PORT string `mapstructure:"PORT"`
	HOST string `mapstructure:"HOST"`
	KEYSPACE string `mapstructure:"KEYSPACE"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			PORT: os.Getenv("PORT"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	// validate config here

	if config.KEYSPACE == "" {
		err = errors.New("KEYSPACE is required")
		return
	}

	return
}

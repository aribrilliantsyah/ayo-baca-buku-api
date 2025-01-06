package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	DB_SOURCE string `mapstructure:"DB_SOURCE"`
}

func LoadAppConfig(path string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

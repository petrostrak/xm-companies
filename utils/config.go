package utils

import (
	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

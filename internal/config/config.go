package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Cache struct {
		Duration int64 `mapstructure:"duration"`
	} `mapstructure:"cache"`
}

func LoadConfig(path string) *Config {
	viper.SetConfigFile(path)
	viper.SetEnvPrefix("OSU") // allows OSU_SERVER_PORT etc.
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &cfg
}

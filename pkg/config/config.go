package config

import (
	"log"

	"github.com/spf13/viper"
)

type App struct {
	Name       string `mapstructure:"name"`
	Version    string `mapstructure:"version"`
	PrivateKey string `mapstructure:"private_key"`
	PublicKey  string `mapstructure:"public_key"`
}

type HTTP struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	App  App  `mapstructure:"app"`
	HTTP HTTP `mapstructure:"http"`
}

// Global config variable
var Cfg *Config

func NewConfig(path string) *Config {
	viper.SetConfigFile(path) // Accept full path (e.g., ./configs/config.yaml)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper.ReadInConfig() has error: %v", err)
	}

	// Optionally override with env variables
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("viper.Unmarshal(&cfg) has error: %v", err)
	}

	Cfg = &cfg

	return Cfg
}

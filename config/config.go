package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN"`
	MessageID     string `envconfig:"MESSAGE_ID"`
}

func LoadConfig() *Config {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err.Error())
	}

	return &config
}

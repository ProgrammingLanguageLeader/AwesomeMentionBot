package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Token                 string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	IsDebugLoggingEnabled bool   `envconfig:"TELEGRAM_IS_DEBUG_LOGGING_ENABLED" default:"false"`
}

func GetConfig() *Config {
	const appPrefix = "amb"
	var config Config
	err := envconfig.Process(appPrefix, &config)
	if err != nil {
		log.Fatal("Cannot process bot config")
	}
	return &config
}

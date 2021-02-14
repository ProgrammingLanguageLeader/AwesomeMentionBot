package setting

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Token                 string `envconfig:"TELEGRAM_TOKEN" required:"true"`
	IsDebugLoggingEnabled bool   `envconfig:"TELEGRAM_IS_DEBUG_LOGGING_ENABLED" default:"false"`
	RedisURL              string `envconfig:"REDIS_URL" default:"localhost:6379"`
	RedisDB               int    `envconfig:"REDIS_DB" default:"0"`
	RedisPassword         string `envconfig:"REDIS_PASSWORD" default:""`
	DevMode               bool   `envconfig:"DEV_MODE" default:"false"`
	BotUsername           string `envconfig:"BOT_USERNAME"`
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

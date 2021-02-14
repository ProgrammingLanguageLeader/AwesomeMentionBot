package db

import (
	"context"
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var client *redis.Client

func GetDBClient() *redis.Client {
	if client == nil {
		config := setting.GetConfig()
		client = redis.NewClient(&redis.Options{
			Addr:     config.RedisURL,
			Password: config.RedisPassword,
			DB:       config.RedisDB,
		})
	}
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logrus.Warn("Redis connection is not stable")
	}
	return client
}

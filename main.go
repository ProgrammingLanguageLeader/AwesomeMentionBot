package main

import (
	"awesomeMentionBot/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := GetConfig()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = config.IsDebugLoggingEnabled
	log.Infof("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// TODO: use webhook for production
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		chat.HandleMessage(bot, &update)
	}
}

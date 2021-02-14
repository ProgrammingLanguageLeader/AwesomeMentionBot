package main

import (
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/chat"
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func main() {
	config := setting.GetConfig()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		logrus.Panic(err)
	}
	bot.Debug = config.IsDebugLoggingEnabled
	logrus.Infof("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// TODO: use webhook for production
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		logrus.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		chat.HandleMessage(bot, &update)
	}
}

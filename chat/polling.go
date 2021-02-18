package chat

import (
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func StartPolling() {
	config := setting.GetConfig()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		logrus.Panic(err)
	}
	bot.Debug = config.IsDebugLoggingEnabled
	logrus.Infof("Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		logrus.Errorf("get updates chan error: %s", err.Error())
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		logrus.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		HandleMessage(bot, &update)
	}
}

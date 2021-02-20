package chat

import (
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"net/http"
)

func StartWebHook() {
	config := setting.GetConfig()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		logrus.Panic(err)
	}
	bot.Debug = config.IsDebugLoggingEnabled
	logrus.Infof("Authorized on account %s", bot.Self.UserName)
	webHookURL := "https://" + config.BotURL + "/" + config.Token
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webHookURL))
	if err != nil {
		logrus.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		logrus.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		logrus.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	go func() {
		err := http.ListenAndServe(":"+config.Port, nil)
		if err != nil {
			logrus.Errorf("webhook error: %v", err)
		}
	}()
	updates := bot.ListenForWebhook("/" + config.Token)
	for update := range updates {
		logrus.Infof("[%s] %s", update.Message.From.UserName, update.Message.Text)
		HandleMessage(bot, &update)
	}
}

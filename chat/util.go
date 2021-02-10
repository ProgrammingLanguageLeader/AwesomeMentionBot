package chat

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func SendMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, responseText string) {
	if responseText == "" {
		return
	}
	response := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
	response.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(response)
	if err != nil {
		log.Errorf("error while sending message: %s", err.Error())
	}
}

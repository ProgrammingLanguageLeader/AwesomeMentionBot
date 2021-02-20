package chat

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strings"
)

func SendMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, responseText string) {
	if responseText == "" {
		return
	}
	response := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
	response.ReplyToMessageID = update.Message.MessageID
	_, err := bot.Send(response)
	if err != nil {
		logrus.Errorf("error while sending message: %v", err)
	}
	logrus.Infof("Message sent: chatID=%d text=%s", response.ChatID, response.Text)
}

func EscapeString(str string) string {
	escapeSeqs := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, escapeSeq := range escapeSeqs {
		escapedSeq := "\\" + escapeSeq
		str = strings.ReplaceAll(str, escapeSeq, escapedSeq)
	}
	return str
}

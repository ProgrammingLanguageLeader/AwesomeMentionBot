package chat

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
)

func HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	var responseText = ""
	if message.Text == "Ещё раз, вы принимаете наркотики?" {
		if rand.Int()%2 == 0 {
			responseText = "https://www.youtube.com/watch?v=iSf2sjKQMo8"
		} else {
			responseText = "https://www.youtube.com/watch?v=AsAxisuXl3o"
		}
		SendMessage(bot, update, responseText)
		return
	}
	skip := message.From.UserName != "dm_shorokhov" &&
		(!message.Chat.IsGroup() || message.Chat.IsGroup() && message.IsCommand())
	if skip && message.IsCommand() {
		responseText = "Не командуй тут мне, я подчиняюсь только хозяину!"
		SendMessage(bot, update, responseText)
		return
	}
	if skip {
		responseText = "Тебе нельзя писать мне! Я на стадии разработки..."
		SendMessage(bot, update, responseText)
		return
	}
	if message.IsCommand() {
		HandleCommand(bot, update)
	}
	SendMessage(bot, update, responseText)
}

package chat

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	var responseText string
	if message.IsCommand() {
		switch message.Command() {
		case "all":
			responseText = "Типа список пользователей"
		case "in":
			responseText = "Лень прикручивать БД"
		case "out":
			responseText = "Очень лень прикручивать БД"
		case "start":
			responseText = "Ну вот ты и поздоровался с этим ботом. И что теперь? Помощь может быть тут /help"
		case "help":
			responseText = "Может быть. А может и нет. А может пошёл ты?"
		default:
			responseText = "Такой команды нет. Помощь может быть тут /help"
		}
	}
	SendMessage(bot, update, responseText)
}

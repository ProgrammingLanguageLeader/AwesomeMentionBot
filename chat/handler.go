package chat

import (
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
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
	config := setting.GetConfig()
	skip := config.DevMode &&
		message.From.UserName != "dm_shorokhov" &&
		(!message.Chat.IsGroup() || message.Chat.IsGroup() && message.IsCommand())
	if skip {
		responseText = "Тебе нельзя писать мне! Я на стадии разработки..."
		SendMessage(bot, update, responseText)
		return
	}
	members := message.NewChatMembers
	wasBotAddedToChat := false
	if members != nil {
		for _, member := range *members {
			if member.UserName == config.BotUsername {
				wasBotAddedToChat = true
			}
		}
	}
	if wasBotAddedToChat || message.ChannelChatCreated || message.GroupChatCreated {
		HandleFirstMessage(bot, update)
	}
	if message.IsCommand() {
		HandleCommand(bot, update)
	}
	SendMessage(bot, update, responseText)
}

package chat

import (
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	config := setting.GetConfig()
	skip := config.DevMode && message.From.UserName != "dm_shorokhov"
	if skip {
		logrus.Debugf("message skipped: from=%s text=%s", message.From.UserName, message.Text)
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
}

package chat

import (
	"container/list"
	"fmt"
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/db"
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	startCommand          = "start"
	allCommand            = "all"
	inCommand             = "in"
	outCommand            = "out"
	helpCommand           = "help"
	setMentionTextCommand = "setmentiontext"
)

const (
	errorMessage = "Something went wrong. Try again later"
	doneMessage  = "Done"
)

func HandleCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	message := update.Message
	if message.IsCommand() {
		switch message.Command() {
		case startCommand:
			HandleFirstMessage(bot, update)
		case allCommand:
			HandleAllCommand(bot, update)
		case inCommand:
			HandleInCommand(bot, update)
		case outCommand:
			HandleOutCommand(bot, update)
		case helpCommand:
			HandleHelpCommand(bot, update)
		case setMentionTextCommand:
			HandleSetMentionText(bot, update)
		default:
			responseText := "No such command. Use /help for getting a manual"
			SendMessage(bot, update, responseText)
		}
	}
}

func HandleFirstMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatUsers := update.Message.NewChatMembers
	config := setting.GetConfig()
	var mentionList []string
	if chatUsers != nil {
		mentionList = make([]string, len(*chatUsers))
		for userIndex, user := range *chatUsers {
			if user.UserName == config.BotUsername {
				continue
			}
			mentionList[userIndex] = user.UserName
		}
	}
	chatID := update.Message.Chat.ID
	settings, _ := db.GetChatSettings(chatID)
	if settings == nil {
		defaultMentionText := "Attention please!"
		db.SaveChatSettings(chatID, &db.ChatSettings{
			MentionText:         defaultMentionText,
			MentionUsernameList: mentionList,
		})
		SendMessage(bot, update, "Bot has been initiated!")
	} else {
		return
	}
}

func HandleAllCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatSettings, err := db.GetChatSettings(update.Message.Chat.ID)
	if err != nil {
		SendMessage(bot, update, "Bot hasn't been initiated. Use /start to do this")
		return
	}
	var replyTextBuilder strings.Builder
	mentionText := chatSettings.MentionText
	commandArgs := update.Message.CommandArguments()
	if commandArgs != "" {
		mentionText = commandArgs
	}
	replyTextBuilder.WriteString(EscapeString(mentionText))
	replyTextBuilder.WriteString("\n")
	for _, username := range chatSettings.MentionUsernameList {
		replyTextBuilder.WriteString(EscapeString(username))
		replyTextBuilder.WriteString(" ")
	}
	for _, user := range chatSettings.MentionUserList {
		replyTextBuilder.WriteString(fmt.Sprintf("[%s](tg://user?id=%d)", user.FirstName, user.ID))
		replyTextBuilder.WriteString(" ")
	}
	response := tgbotapi.NewMessage(update.Message.Chat.ID, replyTextBuilder.String())
	response.ReplyToMessageID = update.Message.MessageID
	response.ParseMode = "MarkdownV2"
	_, err = bot.Send(response)
	if err != nil {
		logrus.Errorf("error while sending message: %v", err)
	}
}

func HandleInCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	messageEntities := *update.Message.Entities
	if len(messageEntities) == 1 {
		SendMessage(bot, update, "Specify one or more username to include them in the mention list")
		return
	}
	includeUsernameList := list.New()
	includeUserList := list.New()
	for _, messageEntity := range messageEntities {
		if messageEntity.Type == "mention" {
			mentionText := text[messageEntity.Offset : messageEntity.Offset+messageEntity.Length]
			includeUsernameList.PushBack(mentionText)
		} else if messageEntity.Type == "text_mention" {
			includeUserList.PushBack(*messageEntity.User)
		}
	}
	_, err := db.IncludeUsersToMentionList(chatID, includeUsernameList, includeUserList)
	if err != nil {
		SendMessage(bot, update, errorMessage)
		return
	}
	SendMessage(bot, update, doneMessage)
}

func HandleOutCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	messageEntities := *update.Message.Entities
	if len(messageEntities) == 1 {
		SendMessage(bot, update, "Specify one or more username to include them in the mention list")
		return
	}
	excludeUsernameList := list.New()
	excludeUserList := list.New()
	for _, messageEntity := range messageEntities {
		if messageEntity.Type == "mention" {
			mentionText := text[messageEntity.Offset : messageEntity.Offset+messageEntity.Length]
			excludeUsernameList.PushBack(mentionText)
		} else if messageEntity.Type == "text_mention" {
			excludeUserList.PushBack(*messageEntity.User)
		}
	}
	_, err := db.ExcludeUsersFromMentionList(chatID, excludeUsernameList, excludeUserList)
	if err != nil {
		SendMessage(bot, update, errorMessage)
		return
	}
	SendMessage(bot, update, doneMessage)
}

func HandleHelpCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	responseText := "My developer is too lazy to write manual..."
	SendMessage(bot, update, responseText)
}

func HandleSetMentionText(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	settings, err := db.GetChatSettings(chatID)
	if err != nil {
		SendMessage(bot, update, errorMessage)
		return
	}
	messageEntities := *update.Message.Entities
	commandEntity := messageEntities[0]
	if len(text) == commandEntity.Length {
		settings.MentionText = ""
	} else {
		trim := strings.Trim(text, " ")
		if len(trim) == commandEntity.Length {
			SendMessage(bot, update, "Incorrect input")
			return
		} else {
			mentionTextOffset := commandEntity.Length + 1
			settings.MentionText = text[mentionTextOffset:]
		}
	}
	db.SaveChatSettings(chatID, settings)
	SendMessage(bot, update, doneMessage)
}

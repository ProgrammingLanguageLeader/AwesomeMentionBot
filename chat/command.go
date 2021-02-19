package chat

import (
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/db"
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
			MentionText: defaultMentionText,
			MentionList: mentionList,
		})
		SendMessage(bot, update, "Bot has been initiated!")
	} else {
		return
	}
}

func HandleAllCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatSettings, _ := db.GetChatSettings(update.Message.Chat.ID)
	var replyTextBuilder strings.Builder
	replyTextBuilder.WriteString(chatSettings.MentionText)
	replyTextBuilder.WriteString("\n")
	for _, username := range chatSettings.MentionList {
		replyTextBuilder.WriteString(username)
		replyTextBuilder.WriteString(" ")
	}
	SendMessage(bot, update, replyTextBuilder.String())
}

func HandleInCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	includeUsernameArray := strings.Split(text, " ")[1:]
	if len(includeUsernameArray) == 0 {
		SendMessage(bot, update, "Specify one of more username to include them in the mention list")
		return
	}
	_, err := db.IncludeUsersToMentionList(chatID, includeUsernameArray)
	if err != nil {
		SendMessage(bot, update, errorMessage)
	}
	SendMessage(bot, update, doneMessage)
}

func HandleOutCommand(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	excludeUsernameArray := strings.Split(text, " ")[1:]
	if len(excludeUsernameArray) == 0 {
		SendMessage(bot, update, "Specify one of more username to include them in the mention list")
		return
	}
	_, err := db.ExcludeUsersFromMentionList(chatID, excludeUsernameArray)
	if err != nil {
		SendMessage(bot, update, errorMessage)
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
	const commandPrefix = "/" + setMentionTextCommand + " "
	settings.MentionText = strings.TrimPrefix(text, commandPrefix)
	db.SaveChatSettings(chatID, settings)
	SendMessage(bot, update, doneMessage)
}

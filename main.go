package main

import (
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/chat"
	"github.com/ProgrammingLanguageLeader/AwesomeMentionBot/setting"
)

func main() {
	config := setting.GetConfig()
	if config.DevMode {
		chat.StartPolling()
	} else {
		chat.StartWebHook()
	}
}

package db

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ChatSettings struct {
	MentionText string   `json:"mention_text"`
	MentionList []string `json:"mention_list"`
}

func SaveChatSettings(chatID int64, settings *ChatSettings) {
	ctx := context.Background()
	settingsMarshal, err := json.Marshal(settings)
	if err != nil {
		logrus.Errorf("marshalling error: %s", err.Error())
	}
	chatKey := strconv.FormatInt(chatID, 10)
	_, err = GetDBClient().Set(ctx, chatKey, settingsMarshal, 0).Result()
	if err != nil {
		logrus.Errorf("create chat settings error: %s", err.Error())
	}
}

func GetChatSettings(chatID int64) (*ChatSettings, error) {
	ctx := context.Background()
	chatKey := strconv.FormatInt(chatID, 10)
	settingsMarshal, err := GetDBClient().Get(ctx, chatKey).Result()
	if err != nil {
		logrus.Warnf("get chat settings error: %s", err.Error())
		return nil, err
	}
	var settings ChatSettings
	if err := json.Unmarshal([]byte(settingsMarshal), &settings); err != nil {
		logrus.Errorf("unmarshal error: %s", err.Error())
	}
	return &settings, nil
}

func IncludeUsersToMentionList(chatID int64, includeUsernameArray []string) (*ChatSettings, error) {
	settings, err := GetChatSettings(chatID)
	if err != nil {
		logrus.Errorf("adding user error: %s", err.Error())
		return nil, err
	}

	mentionUsernameSet := map[string]bool{}
	needSave := false
	for _, username := range settings.MentionList {
		mentionUsernameSet[username] = true
	}
	for _, includeUsername := range includeUsernameArray {
		_, usernameInList := mentionUsernameSet[includeUsername]
		if !usernameInList {
			mentionUsernameSet[includeUsername] = true
			needSave = true
		}
	}

	if needSave {
		newUsernameArray := make([]string, len(mentionUsernameSet))
		for username := range mentionUsernameSet {
			newUsernameArray = append(newUsernameArray, username)
		}
		settings.MentionList = newUsernameArray
		SaveChatSettings(chatID, settings)
	}
	return settings, nil
}

func ExcludeUsersFromMentionList(chatID int64, excludeUsernameArray []string) (*ChatSettings, error) {
	settings, err := GetChatSettings(chatID)
	if err != nil {
		logrus.Errorf("exclude user error: %s", err.Error())
		return nil, err
	}

	mentionUsernameSet := map[string]bool{}
	needSave := false
	for _, username := range settings.MentionList {
		mentionUsernameSet[username] = true
	}
	for _, excludeUsername := range excludeUsernameArray {
		_, usernameInList := mentionUsernameSet[excludeUsername]
		if usernameInList {
			delete(mentionUsernameSet, excludeUsername)
			needSave = true
		}
	}

	if needSave {
		newUsernameArray := make([]string, len(mentionUsernameSet))
		for username := range mentionUsernameSet {
			newUsernameArray = append(newUsernameArray, username)
		}
		settings.MentionList = newUsernameArray
		SaveChatSettings(chatID, settings)
	}
	return settings, nil
}

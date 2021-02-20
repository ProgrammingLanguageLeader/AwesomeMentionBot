package db

import (
	"container/list"
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ChatSettings struct {
	MentionText         string          `json:"mention_text"`
	MentionUsernameList []string        `json:"mention_username_list"`
	MentionUserList     []tgbotapi.User `json:"mention_user_list"`
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

func IncludeUsersToMentionList(chatID int64, includeUsernameList *list.List, includeUserList *list.List) (*ChatSettings, error) {
	settings, err := GetChatSettings(chatID)
	if err != nil {
		logrus.Errorf("adding user error: %s", err.Error())
		return nil, err
	}

	mentionUsernameSet := map[string]bool{}
	needSave := false
	for _, username := range settings.MentionUsernameList {
		mentionUsernameSet[username] = true
	}
	mentionUserSet := map[int]tgbotapi.User{}
	for _, user := range settings.MentionUserList {
		mentionUserSet[user.ID] = user
	}

	for includeUsername := includeUsernameList.Front(); includeUsername != nil; includeUsername = includeUsername.Next() {
		username := includeUsername.Value.(string)
		_, usernameInList := mentionUsernameSet[username]
		if !usernameInList {
			mentionUsernameSet[username] = true
			needSave = true
		}
	}
	for includeUser := includeUserList.Front(); includeUser != nil; includeUser = includeUser.Next() {
		user := includeUser.Value.(tgbotapi.User)
		_, userInList := mentionUserSet[user.ID]
		if !userInList {
			mentionUserSet[user.ID] = user
			needSave = true
		}
	}

	if needSave {
		newUsernameArray := make([]string, len(mentionUsernameSet))
		usernameIndex := 0
		for username := range mentionUsernameSet {
			newUsernameArray[usernameIndex] = username
			usernameIndex++
		}
		settings.MentionUsernameList = newUsernameArray
		newUserArray := make([]tgbotapi.User, len(mentionUserSet))
		userIndex := 0
		for _, user := range mentionUserSet {
			newUserArray[userIndex] = user
			userIndex++
		}
		settings.MentionUserList = newUserArray
		SaveChatSettings(chatID, settings)
	}
	return settings, nil
}

func ExcludeUsersFromMentionList(chatID int64, excludeUsernameList *list.List, excludeUserList *list.List) (*ChatSettings, error) {
	settings, err := GetChatSettings(chatID)
	if err != nil {
		logrus.Errorf("exclude user error: %s", err.Error())
		return nil, err
	}

	mentionUsernameSet := map[string]bool{}
	needSave := false
	for _, username := range settings.MentionUsernameList {
		mentionUsernameSet[username] = true
	}
	mentionUserSet := map[int]tgbotapi.User{}
	for _, user := range settings.MentionUserList {
		mentionUserSet[user.ID] = user
	}

	for excludeUsername := excludeUsernameList.Front(); excludeUsername != nil; excludeUsername = excludeUsername.Next() {
		username := excludeUsername.Value.(string)
		_, usernameInList := mentionUsernameSet[username]
		if usernameInList {
			delete(mentionUsernameSet, username)
			needSave = true
		}
	}
	for excludeUser := excludeUserList.Front(); excludeUser != nil; excludeUser = excludeUser.Next() {
		user := excludeUser.Value.(tgbotapi.User)
		_, userInList := mentionUserSet[user.ID]
		if userInList {
			delete(mentionUserSet, user.ID)
			needSave = true
		}
	}

	if needSave {
		newUsernameArray := make([]string, len(mentionUsernameSet))
		usernameIndex := 0
		for username := range mentionUsernameSet {
			newUsernameArray[usernameIndex] = username
			usernameIndex++
		}
		settings.MentionUsernameList = newUsernameArray
		newUserArray := make([]tgbotapi.User, len(mentionUserSet))
		userIndex := 0
		for _, user := range mentionUserSet {
			newUserArray[userIndex] = user
			userIndex++
		}
		settings.MentionUserList = newUserArray
		SaveChatSettings(chatID, settings)
	}
	return settings, nil
}

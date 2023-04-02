package model

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// RerenderMenu - is very specific case functionality.
// Re-renders currently opened menu (with MessageId) in chat with ChatId
// with one localized (if such exists) of language in Ctx with key menu.LocaleContextKey
type RerenderMenu struct {
	ChatId    int64
	MessageId int64
	Ctx       context.Context
}

func (u *RerenderMenu) GetMessage() tgbotapi.Chattable {
	return nil
}

func (u *RerenderMenu) GetChatId() int64 {
	return u.ChatId
}

func (u *RerenderMenu) SetChatIdIfZero(c int64) {
	if u.ChatId == 0 {
		u.ChatId = c
	}
}

func (u *RerenderMenu) ForceSetMessageId(id int64) {
	u.MessageId = id
}

func (u *RerenderMenu) GetMessageId() int64 {
	return u.MessageId
}

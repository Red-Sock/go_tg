package model

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OpenMenu struct {
	ChatId    int64
	MessageId int64
	Command   string
	Ctx       context.Context
}

func NewOpenMenu(ctx context.Context, menuCommand string) *OpenMenu {
	return &OpenMenu{
		Ctx:     ctx,
		Command: menuCommand,
	}
}

func (u *OpenMenu) GetMessage() tgbotapi.Chattable {
	return nil
}

func (u *OpenMenu) GetChatId() int64 {
	return u.ChatId
}

func (u *OpenMenu) SetChatIdIfZero(c int64) {
	if u.ChatId == 0 {
		u.ChatId = c
	}
}

func (u *OpenMenu) SetMessageId(id int64) {
	u.MessageId = id
}

func (u *OpenMenu) GetMessageId() int64 {
	return u.MessageId
}

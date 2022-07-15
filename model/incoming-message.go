package model

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageIn struct {
	Command string
	Args    []string
	Ctx     context.Context

	*tgbotapi.Message
}

func (m *MessageIn) Delete() {
	//TODO
}

func (m *MessageIn) Response(out string) {
	//internal.Send(NewReplyToMessage(out, m.Chat.ID, int64(m.MessageID)))
}

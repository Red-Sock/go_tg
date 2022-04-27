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
	ChatMsg
}

type ChatMsg struct {
}

func (m *ChatMsg) Delete() {
	// TODO
}

func (m *ChatMsg) Response(out MessageOut) {

}

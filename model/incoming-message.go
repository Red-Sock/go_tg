package model

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageIn struct {
	Command    string
	Args       []string
	Ctx        context.Context
	IsCallback bool
	*tgbotapi.Message
}

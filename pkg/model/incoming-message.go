package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageIn struct {
	*tgbotapi.Message

	Command string
	Args    []string
}

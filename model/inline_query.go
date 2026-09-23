package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// InlineQueryIn wraps an incoming Telegram inline query - sent when a user
// types "@botname <query>" in any chat. Unlike MessageIn, there is no
// associated chat to reply into: the answer goes back via answerInlineQuery,
// not through a Chat.
type InlineQueryIn struct {
	*tgbotapi.InlineQuery
}

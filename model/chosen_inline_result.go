package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ChosenInlineResultIn wraps an incoming Telegram chosen inline result - sent
// when a user actually taps one of the results answered for an inline query.
// Unlike InlineQueryIn, there is no answer to send back: this is a
// notification only, typically used to lazily finish preparing the chosen
// result's content (e.g. via InlineMessageID).
type ChosenInlineResultIn struct {
	*tgbotapi.ChosenInlineResult
}

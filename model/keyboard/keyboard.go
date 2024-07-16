package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboard interface {
	ToInlineMarkup() *tgbotapi.InlineKeyboardMarkup
	ToKeyboardMarkup() *tgbotapi.ReplyKeyboardMarkup

	IsReplyKeyboard() bool
}

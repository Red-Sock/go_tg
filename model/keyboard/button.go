package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Button struct {
	Text         string
	Value        string
	IsStandAlone bool

	InternalButton               *tgbotapi.KeyboardButton
	InternalInlineKeyboardButton *tgbotapi.InlineKeyboardButton
}

func NewButton(text, value string) Button {
	return Button{
		Text:         text,
		Value:        value,
		IsStandAlone: false,
	}
}

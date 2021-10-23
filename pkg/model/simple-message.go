package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageOut struct {
	Text          string
	InlineButtons *InlineKeyboard
	Keyboard      *Keyboard
}

func (m *MessageOut) ToTgMessage(chatId int64) tgbotapi.Chattable {
	message := tgbotapi.NewMessage(chatId, m.Text)

	if m.InlineButtons != nil {
		message.BaseChat.ReplyMarkup = m.InlineButtons.ToMarkup()
	}
	if m.Keyboard != nil {
		message.BaseChat.ReplyMarkup = m.Keyboard.toMarkup()
	}

	return &message
}

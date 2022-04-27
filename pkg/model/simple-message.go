package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageOut struct {
	MessageId     int
	ChatId        int64
	Text          string
	InlineButtons *InlineKeyboard
	Keyboard      *Keyboard
}

func (m *MessageOut) Execute(bot *tgbotapi.BotAPI) {
	m.Send(bot, m.ChatId)
}

func (m *MessageOut) Send(bot *tgbotapi.BotAPI, chatId int64) error {
	message := tgbotapi.NewMessage(chatId, m.Text)

	if m.InlineButtons != nil {
		message.BaseChat.ReplyMarkup = m.InlineButtons.ToMarkup()
	}
	if m.Keyboard != nil {
		message.BaseChat.ReplyMarkup = m.Keyboard.toMarkup()
	}
	message.DisableWebPagePreview = true

	msg, err := bot.Send(message)
	if err != nil {
		return err
	}

	m.MessageId = msg.MessageID
	return nil
}

func (m *MessageOut) SetChatIdIfZero(i int64) {
	if m.ChatId == 0 {
		m.ChatId = i
	}
}

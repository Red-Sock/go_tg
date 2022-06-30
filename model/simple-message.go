package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageOut struct {
	MessageId int64
	ChatId    int64
	Text      string
}

func (m *MessageOut) GetMessage() tgbotapi.Chattable {
	message := tgbotapi.NewMessage(m.ChatId, m.Text)
	message.DisableWebPagePreview = true

	return message
}

func (m *MessageOut) SetChatIdIfZero(i int64) {
	if m.ChatId == 0 {
		m.ChatId = i
	}
}

func (m *MessageOut) GetChatId() int64 {
	return m.ChatId
}

func (m *MessageOut) SetMessageId(id int64) {
	m.MessageId = id
}

func (m *MessageOut) GetMessageId() int64 {
	return m.MessageId
}

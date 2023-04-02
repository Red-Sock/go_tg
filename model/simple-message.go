package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageOut interface {
	//GetText() string
	SetChatIdIfZero(chatID int64)
	GetChatId() int64

	ForceSetMessageId(msgId int64)

	GetMessage() tgbotapi.Chattable
}

type messageOut struct {
	MessageId int64
	ChatId    int64
	Text      string
}

func NewMessage(text string) *messageOut {
	return &messageOut{Text: text}
}

func NewMessageToChat(text string, chatId int64) *messageOut {
	return &messageOut{Text: text, ChatId: chatId}
}

func NewReplyToMessage(msg string, chatId int64, messageId int64) *messageOut {
	return &messageOut{
		MessageId: messageId,
		ChatId:    chatId,
		Text:      msg,
	}
}

func (m *messageOut) GetMessage() tgbotapi.Chattable {
	message := tgbotapi.NewMessage(m.ChatId, m.Text)
	message.DisableWebPagePreview = true

	return message
}

func (m *messageOut) SetChatIdIfZero(i int64) {
	if m.ChatId == 0 {
		m.ChatId = i
	}
}

func (m *messageOut) GetChatId() int64 {
	return m.ChatId
}

func (m *messageOut) ForceSetMessageId(id int64) {
	m.MessageId = id
}

func (m *messageOut) GetMessageId() int64 {
	return m.MessageId
}

func (m *messageOut) GetText() string {
	return m.Text
}

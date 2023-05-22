package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DeleteMessage struct {
	ChatId    int64
	MessageId int64
}

func (m *DeleteMessage) GetMessage() tgbotapi.Chattable {
	return tgbotapi.NewDeleteMessage(m.ChatId, int(m.MessageId))
}

func (m *DeleteMessage) SetChatIdIfZero(chatID int64) {
	m.ChatId = chatID
}

func (m *DeleteMessage) GetChatId() int64 {
	return m.ChatId
}

func (m *DeleteMessage) ForceSetMessageId(msgId int64) {
	m.MessageId = msgId
}

func (m *DeleteMessage) GetMessageId() int64 {
	return m.MessageId
}

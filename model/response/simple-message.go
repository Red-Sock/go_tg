package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Red-Sock/go_tg/model/keyboard"
)

type MessageOut struct {
	ChatId int64
	Text   string

	MessageId int64

	Keys     *keyboard.InlineKeyboard
	Entities []tgbotapi.MessageEntity
}

func NewMessage(text string) *MessageOut {
	return &MessageOut{Text: text}
}

func NewMessageToChat(text string, chatId int64) *MessageOut {
	return &MessageOut{Text: text, ChatId: chatId}
}

func (m *MessageOut) GetMessage() tgbotapi.Chattable {
	if m.MessageId != 0 {
		msg := EditMessage{
			ChatId:    m.ChatId,
			Text:      m.Text,
			MessageId: m.MessageId,
			Keys:      m.Keys,
		}
		return msg.GetMessage()
	}

	message := tgbotapi.NewMessage(m.ChatId, m.Text)
	message.DisableWebPagePreview = true

	if len(m.Entities) != 0 {
		message.Entities = m.Entities
	}

	if m.Keys != nil {
		keyboard := m.Keys.ToMarkup()
		message.ReplyMarkup = keyboard
	}

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

func (m *MessageOut) ForceSetMessageId(id int64) {
	m.MessageId = id
}

func (m *MessageOut) GetMessageId() int64 {
	return m.MessageId
}

func (m *MessageOut) AddKeyboard(keys keyboard.InlineKeyboard) {
	m.Keys = &keys
}

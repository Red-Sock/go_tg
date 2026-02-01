package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Red-Sock/go_tg/model/keyboard"
)

type EditMessage struct {
	ChatId int64
	Text   string

	MessageId int64

	Keys     keyboard.Keyboard
	Entities []tgbotapi.MessageEntity
}

func (m *EditMessage) GetMessage() tgbotapi.Chattable {
	switch {
	case m.Keys != nil && m.Text != "":
		mu := m.Keys.ToInlineMarkup()
		message := tgbotapi.NewEditMessageTextAndMarkup(m.ChatId, int(m.MessageId), m.Text, *mu)
		message.DisableWebPagePreview = true
		message.Entities = m.Entities
		return message
	case m.Keys != nil:
		mu := m.Keys.ToInlineMarkup()
		message := tgbotapi.NewEditMessageReplyMarkup(m.ChatId, int(m.MessageId), *mu)

		return message
	case m.Text != "":
		message := tgbotapi.NewEditMessageText(m.ChatId, int(m.MessageId), m.Text)
		message.Entities = m.Entities
		return message
	default:
		return nil
	}
}

func (m *EditMessage) SetChatIdIfZero(chatID int64) {
	m.ChatId = chatID
}

func (m *EditMessage) GetChatId() int64 {
	return m.ChatId
}

func (m *EditMessage) ForceSetMessageId(msgId int64) {
	m.MessageId = msgId
}

func (m *EditMessage) GetMessageId() int64 {
	return m.MessageId
}

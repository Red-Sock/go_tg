package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Red-Sock/go_tg/model/keyboard"
	"github.com/Red-Sock/go_tg/model/media"
)

type MessageOut struct {
	ChatId int64
	Text   string

	MessageId int64

	Keys           keyboard.Keyboard
	Entities       []tgbotapi.MessageEntity
	ReplyMessageId int64

	Media []media.Media
}

func NewMessage(text string, opts ...opt) *MessageOut {
	m := &MessageOut{Text: text}

	for _, o := range opts {
		o(m)
	}

	return m
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

	if len(m.Media) != 0 {
		if len(m.Media) == 1 {
			return m.asSingleMedia()
		}
		return m.asMediaGroup()
	}

	message := tgbotapi.NewMessage(m.ChatId, m.Text)
	message.DisableWebPagePreview = true

	if len(m.Entities) != 0 {
		message.Entities = m.Entities
	}

	if m.Keys != nil {
		var keyboard any
		if m.Keys.IsReplyKeyboard() {
			keyboard = m.Keys.ToKeyboardMarkup()
		} else {
			keyboard = m.Keys.ToInlineMarkup()
		}

		message.ReplyMarkup = keyboard
	}

	message.ReplyToMessageID = int(m.ReplyMessageId)

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

func (m *MessageOut) AddKeyboard(keys keyboard.Keyboard) {
	m.Keys = keys
}

func (m *MessageOut) asSingleMedia() tgbotapi.Chattable {
	return m.Media[0].AsSingleTgMedia(m.ChatId)
}

func (m *MessageOut) asMediaGroup() tgbotapi.Chattable {
	var files []any

	for _, file := range m.Media {
		files = append(files, file.AsInputMedia())
	}

	return tgbotapi.NewMediaGroup(m.ChatId, files)
}

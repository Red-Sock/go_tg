package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageIn struct {
	*tgbotapi.Message

	Command string
	Args    []string
}

type MessageOut struct {
	Text          string
	InlineButtons *InlineKeyboard
	Keyboard      *Keyboard
}

func (m *MessageOut) ToTgMessage(chatId int64) *tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(chatId, m.Text)

	if m.InlineButtons != nil {
		message.BaseChat.ReplyMarkup = m.InlineButtons.ToMarkup()
	}
	if m.Keyboard != nil {
		message.BaseChat.ReplyMarkup = m.Keyboard.toMarkup()
	}

	return &message
}

type MessageChan chan Conversation

type Conversation struct {
	MessageOut *MessageOut
	MessageIn  *MessageIn
}

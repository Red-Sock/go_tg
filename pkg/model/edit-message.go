package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageEdit struct {
	MessageId int
	ChatId    int64
	Text      *string

	ReplyKeyboard *InlineKeyboard
}

func EditMessageText(text string, messageId int) *MessageEdit {
	return &MessageEdit{
		MessageId: messageId,
		Text:      &text,
	}
}

func EditMessageReply(keyboard *InlineKeyboard, messageId int) *MessageEdit {
	return &MessageEdit{
		MessageId:     messageId,
		ReplyKeyboard: keyboard,
	}
}

func (m *MessageEdit) Send(bot *tgbotapi.BotAPI, chatId int64) (err error) {
	var msg tgbotapi.Message
	if m.Text != nil {
		msg, err = bot.Send(tgbotapi.NewEditMessageText(chatId, m.MessageId, *m.Text))
	}

	if m.ReplyKeyboard != nil {
		msg, err = bot.Send(tgbotapi.NewEditMessageReplyMarkup(chatId, m.MessageId, *m.ReplyKeyboard.ToMarkup()))
	}

	m.MessageId = msg.MessageID

	return nil
}

func (m *MessageEdit) Execute(bot *tgbotapi.BotAPI) {
	m.Send(bot, m.ChatId)
}

func (m *MessageEdit) SetChatIdIfZero(i int64) {
	if m.ChatId == 0 {
		m.ChatId = i
	}
}

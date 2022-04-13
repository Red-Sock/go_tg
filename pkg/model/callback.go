package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackType int

const (
	Callback_Type_Bad CallbackType = iota
	Callback_Type_CallCommand
	Callback_Type_OpenMenu
	Callback_Type_TransitToMenu
)

type Callback struct {
	Command string
	Args    []string
	Type    CallbackType
	Menu    Menu

	Text string

	ReplyMarkup *tgbotapi.InlineKeyboardMarkup
}

// Send sends a message back to user if nessesery
func (c *Callback) Send(api *tgbotapi.BotAPI, chatId int64) (err error) {
	var mc tgbotapi.MessageConfig

	if len(c.Text) != 0 {
		mc = tgbotapi.NewMessage(chatId, c.Text)
	}

	mc.ReplyMarkup = c.ReplyMarkup

	if len(mc.Text) != 0 {
		_, err = api.Send(mc)
	}
	return err
}

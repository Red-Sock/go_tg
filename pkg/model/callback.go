package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackType int

const (
	Callback_Type_OpenMenu = iota
	Callback_Type_TransitToMenu
)

type Callback struct {
	Command string
	Type    CallbackType
}

func (c *Callback) Send(api *tgbotapi.BotAPI, chatId int64) error {
	return nil
}

func (c *Callback) Callback(message *tgbotapi.Message) {
	switch c.Type {
	case Callback_Type_OpenMenu:
		message.Text = OpenMenu + " " + c.Command
	case Callback_Type_TransitToMenu:
		message.Text = MenuCall + " " + c.Command
	}
}

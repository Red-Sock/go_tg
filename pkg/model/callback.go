package model

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackType int

const (
	Callback_Type_CallCommand CallbackType = iota
	Callback_Type_OpenMenu
	Callback_Type_TransitToMenu
)

type Callback struct {
	Command string
	Type    CallbackType
	Menu    Menu

	Text string

	Ctx         context.Context
	replyMarkup *tgbotapi.InlineKeyboardMarkup
}

// Send sends a message back to user if nessesery
func (c *Callback) Send(api *tgbotapi.BotAPI, chatId int64) (err error) {
	var mc tgbotapi.MessageConfig

	if len(c.Text) != 0 {
		mc = tgbotapi.NewMessage(chatId, c.Text)
	}

	mc.ReplyMarkup = c.replyMarkup

	if len(mc.Text) != 0 {
		_, err = api.Send(mc)
	}
	return err
}

// Callback modifies tgbotapi.Message message so it can be rehandled
// send defines wether message should be send instantly
// or rehandled
func (c *Callback) Process(message *MessageIn) (send bool) {
	switch c.Type {
	case Callback_Type_CallCommand:
		message.Text = c.Command
	case Callback_Type_OpenMenu:
		if c.Menu != nil {
			menu := c.Menu.GetPage(c.Ctx, 0)
			c.replyMarkup = menu.ToMarkup()
			return true
		}
		message.Text = OpenMenu + " " + c.Command
	case Callback_Type_TransitToMenu:
		message.Text = MenuCall + " " + c.Command
	}
	return false
}

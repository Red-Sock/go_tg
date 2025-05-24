package internal

import (
	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
)

var d interfaces.CommandHandler = &DefaultHandler{}

type DefaultHandler struct {
	Logger logrus.FieldLogger
}

func (d *DefaultHandler) Handle(in *model.MessageIn, out interfaces.Chat) error {
	msg := "Couldn't handle \"" + in.Text + "\" command"
	d.Logger.Error(msg)
	return out.SendMessage(response.NewMessage(msg))
}

func (d *DefaultHandler) GetDescription() string {
	return ""
}

func (d *DefaultHandler) GetCommand() string {
	return ""
}

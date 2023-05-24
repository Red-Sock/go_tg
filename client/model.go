package client

import (
	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
)

type chatHandler struct {
	handler interfaces.CommandHandler
	msgCh   chan *model.MessageIn
}

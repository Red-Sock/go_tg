package client

import "github.com/Red-Sock/go_tg/interfaces"

type chatHandler struct {
	handler interfaces.CommandHandler
	msgCh   chan *model.MessageIn
}

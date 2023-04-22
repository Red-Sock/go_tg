package client

import "github.com/RedSock/go_tg/interfaces"

type chatHandler struct {
	command string
	handler interfaces.CommandHandler
}

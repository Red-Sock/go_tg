package client

import "github.com/Red-Sock/go_tg/interfaces"

type chatHandler struct {
	command string
	handler interfaces.CommandHandler
}

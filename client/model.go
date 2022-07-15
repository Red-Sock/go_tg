package client

import "github.com/AlexSkilled/go_tg/interfaces"

type chatHandler struct {
	command string
	handler interfaces.CommandHandler
}

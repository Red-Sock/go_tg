package client

import (
	"github.com/AlexSkilled/go_tg/interfaces"
	"github.com/AlexSkilled/go_tg/model"
)

type chatHandler struct {
	handler interfaces.CommandHandler
	msgCh   chan *model.MessageIn
}

package internal

import (
	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
)

type ChatHandler struct {
	Handler interfaces.CommandHandler
	MsgCh   chan *model.MessageIn
}

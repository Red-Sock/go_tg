package internal

import (
	"time"

	"github.com/Red-Sock/go_tg/interfaces"
)

type Chat struct {
	ChatId  int64
	COut    func(out interfaces.MessageOut) error
	Timeout time.Duration
}

func (r *Chat) SendMessage(ins interfaces.MessageOut) error {
	ins.SetChatIdIfZero(r.ChatId)
	return r.COut(ins)
}

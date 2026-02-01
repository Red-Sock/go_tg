package internal

import (
	"time"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
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

func (r *Chat) DeleteIncomingMessage(in *model.MessageIn) {
	deleteMsg := &response.DeleteMessage{
		ChatId:    in.Chat.ID,
		MessageId: int64(in.MessageID),
	}

	_ = r.COut(deleteMsg)
}

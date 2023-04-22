package client

import (
	"context"
	"time"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
)

type chat struct {
	chatId  int64
	cOut    chan<- interfaces.MessageOut
	cIn     <-chan *model.MessageIn
	timeout time.Duration
}

func (r *chat) SendMessage(ins interfaces.MessageOut) {
	ins.SetChatIdIfZero(r.chatId)
	r.cOut <- ins
}

func (r *chat) GetInput(ctx context.Context) (*model.MessageIn, error) {

	_, ok := ctx.Deadline()
	if !ok {
		ctx, _ = context.WithTimeout(ctx, r.timeout)
	}
	select {
	case <-ctx.Done():
		return nil, interfaces.ErrTimeout
	case msg := <-r.cIn:
		return msg, nil
	}
}

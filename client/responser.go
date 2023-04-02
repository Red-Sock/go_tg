package client

import (
	"context"
	"time"

	"github.com/AlexSkilled/go_tg/interfaces"
	"github.com/AlexSkilled/go_tg/model"
)

type chat struct {
	chatId int64
	cOut   chan<- model.MessageOut
	cIn    <-chan *model.MessageIn
}

func (r *chat) SendMessage(ins model.MessageOut) {
	ins.SetChatIdIfZero(r.chatId)
	r.cOut <- ins
}

func (r *chat) GetInput(ctx context.Context) (*model.MessageIn, error) {
	ctx, _ = context.WithTimeout(ctx, time.Second*10)
	select {
	case <-ctx.Done():
		return nil, interfaces.ErrTimeout
	case msg := <-r.cIn:
		return msg, nil
	}
}

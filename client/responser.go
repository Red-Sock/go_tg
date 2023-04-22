package client

import "github.com/Red-Sock/go_tg/interfaces"

type responser struct {
	chatId int64
	c      chan<- interfaces.Instruction
}

func (r *responser) Send(ins interfaces.Instruction) {
	ins.SetChatIdIfZero(r.chatId)
	r.c <- ins
}

package client

import "github.com/AlexSkilled/go_tg/interfaces"

type Responser struct {
	chatId int64
	c      chan<- interfaces.Instruction
}

func (r *Responser) Send(ins interfaces.Instruction) {
	ins.SetChatIdIfZero(r.chatId)
	r.c <- ins
}

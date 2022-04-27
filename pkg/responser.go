package pkg

type Responser struct {
	chatId int64
	c      chan<- Instruction
}

func (r *Responser) Send(ins Instruction) {
	ins.SetChatIdIfZero(r.chatId)
	r.c <- ins
}

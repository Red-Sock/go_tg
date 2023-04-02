package internal

import (
	"github.com/AlexSkilled/go_tg/interfaces"
	"github.com/AlexSkilled/go_tg/model"
)

var s chan<- model.MessageOut

func SetSender(sender chan<- model.MessageOut) {
	s = sender
}

func Send(instruction interfaces.Instruction) {
	s <- instruction
}

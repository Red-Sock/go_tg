package internal

import (
	"github.com/AlexSkilled/go_tg/interfaces"
)

var s chan<- interfaces.MessageOut

func SetSender(sender chan<- interfaces.MessageOut) {
	s = sender
}

func Send(instruction interfaces.MessageOut) {
	s <- instruction
}

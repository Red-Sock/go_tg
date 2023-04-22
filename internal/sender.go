package internal

import "github.com/Red-Sock/go_tg/interfaces"

var s chan<- interfaces.Instruction

func SetSender(sender chan<- interfaces.Instruction) {
	s = sender
}

func Send(instruction interfaces.Instruction) {
	s <- instruction
}

package internal

import "github.com/AlexSkilled/go_tg/interfaces"

var s chan<- interfaces.Instruction

func SetSender(sender chan<- interfaces.Instruction) {
	s = sender
}

func Send(instruction interfaces.Instruction) {
	s <- instruction
}

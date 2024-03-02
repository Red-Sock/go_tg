package send

import "github.com/Red-Sock/go_tg/interfaces"

var send func(out interfaces.MessageOut) error

func SetSender(s func(out interfaces.MessageOut) error) {
	send = s
}

func Send(instruction interfaces.MessageOut) {
	send(instruction)
}

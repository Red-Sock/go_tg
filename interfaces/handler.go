package interfaces

import (
	"github.com/Red-Sock/go_tg/model"
)

type Handler interface {
	// Handle function that gets executed when user send certain command
	// Use it When you have to work with dynamic data, such as
	// database, other service and etc.
	// in - is the message with payload in context.
	// out - is the channel which will send outgoing messages
	Handle(in *model.MessageIn, out Chat) error
}

// CommandHandler - structure that contains logic for
// handling commands that starts with "/" symbol
type CommandHandler interface {
	Handler
	GetCommand() string
}

type Description interface {
	GetDescription() string
}

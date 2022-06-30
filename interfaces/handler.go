package interfaces

import "github.com/AlexSkilled/go_tg/model"

// CommandHandler - structure that contains logic for
// handling commands that starts with "/" symbol
type CommandHandler interface {
	// Handle function that gets executed when user send certain command
	// Use it When you have to work with dynamic data, such as
	// database, other service and etc.
	// in - is the message with payload in context.
	// out - is the channel which will send outgoing messages
	Handle(in *model.MessageIn, out Sender)
	// Dump is has to clear all context of current handler with
	// chatId
	Dump(chatId int64)
}

package interfaces

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Sender TODO передалать на менее общий интерфейс
type Sender interface {
	Send(Instruction)
}

type Instruction interface {
	GetMessage() tgbotapi.Chattable

	GetChatId() int64
	SetChatIdIfZero(c int64)

	SetMessageId(id int64)
	GetMessageId() int64
}

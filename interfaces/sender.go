package interfaces

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

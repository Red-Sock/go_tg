package interfaces

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/AlexSkilled/go_tg/model"
)

var ErrTimeout = errors.New("timeout waiting response from user")

type Chat interface {
	GetInput(ctx context.Context) (*model.MessageIn, error)
	SendMessage(out model.MessageOut)
}

// Sender TODO передалать на менее общий интерфейс
type Sender interface {
	Send(Instruction)
}

type Instruction interface {
	GetMessage() tgbotapi.Chattable

	GetChatId() int64
	SetChatIdIfZero(c int64)

	ForceSetMessageId(id int64)
	GetMessageId() int64
}

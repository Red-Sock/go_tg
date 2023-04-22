package interfaces

import (
	"context"
	"errors"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Red-Sock/go_tg/model"
)

var ErrTimeout = errors.New("timeout waiting response from user")

const UserResponseTimeout = time.Second * 10

type Chat interface {
	// GetInput - awaits for users response with given in a ctx timeout or
	GetInput(ctx context.Context) (*model.MessageIn, error)
	SendMessage(out MessageOut)
}

type MessageOut interface {
	SetChatIdIfZero(chatID int64)
	GetChatId() int64

	ForceSetMessageId(msgId int64)
	GetMessageId() int64

	GetMessage() tgbotapi.Chattable
}

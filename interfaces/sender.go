package interfaces

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const UserResponseTimeout = time.Second * 10

type Chat interface {
	SendMessage(out MessageOut) error
}

type MessageOut interface {
	SetChatIdIfZero(chatID int64)
	GetChatId() int64

	ForceSetMessageId(msgId int64)
	GetMessageId() int64

	GetMessage() tgbotapi.Chattable
}

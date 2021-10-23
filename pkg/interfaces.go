package pkg

import (
	"context"

	"github.com/AlexSkilled/go_tg/pkg/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CommandHandler - structure that contains logic for
// handling commands that starts with "/" symbol
// Handle - function that gets executed when user send certain command
// Use it When you have to work with dynamic data, such as
// database, other service and etc
type CommandHandler interface {
	Handle(ctx context.Context, in *model.MessageIn) (out *model.MessageOut)
}

// CommandHandlerFunc - function for handling simple request logic
// use it when you work with static data and
// don't need to create a structure
type CommandHandlerFunc func(ctx context.Context, in *model.MessageIn) (out *model.MessageOut)

func (c CommandHandlerFunc) Handle(ctx context.Context, in *model.MessageIn) (out *model.MessageOut) {
	return c(ctx, in)
}

type EnrichContext interface {
	GetContext(in *model.MessageIn) (context.Context, error)
}
type GetContextFunc func(in *model.MessageIn) (context.Context, error)

func (g GetContextFunc) GetContext(in *model.MessageIn) (context.Context, error) {
	return g(in)
}

type TgMessage interface {
	ToTgMessage(chatId int64) tgbotapi.Chattable
}

package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlexSkilled/go_tg/pkg/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CommandHandler - structure that contains logic for
// handling commands that starts with "/" symbol
// Handle - function that gets executed when user send certain command
// Use it When you have to work with dynamic data, such as
// database, other service and etc
type CommandHandler interface {
	Handle(in *model.MessageIn, out Sender)
	Dump(id int64)
}

type ExternalContext interface {
	GetContext(in *model.MessageIn) (context.Context, error)
}
type GetContextFunc func(in *model.MessageIn) (context.Context, error)

func (g GetContextFunc) GetContext(in *model.MessageIn) (context.Context, error) {
	return g(in)
}

type TgMessage interface {
	Send(api *tgbotapi.BotAPI, chatId int64) error
}

type MultipleMessage map[int64][]TgMessage

func (m *MultipleMessage) Send(bot *tgbotapi.BotAPI, _ int64) error {
	errors := make([]string, 0)
	for chatId, messages := range *m {
		for _, message := range messages {
			err := message.Send(bot, chatId)
			if err != nil {
				errors = append(errors, err.Error())
			}
		}

	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("%s", strings.Join(errors, "\n"))
}

type Sender interface {
	Send(Instruction)
}

type Instruction interface {
	Execute(bot *tgbotapi.BotAPI)
	SetChatIdIfZero(c int64)
}

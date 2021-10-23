package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlexSkilled/go_tg/pkg/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// Bot - allows you to interact with telegram bot
// with some features
type Bot struct {
	Bot      *tgbotapi.BotAPI
	chats    map[int64]CommandHandler
	handlers map[string]CommandHandler
	EnrichContext
	separator string
}

func NewBot(token string) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &Bot{
		Bot:       bot,
		chats:     make(map[int64]CommandHandler),
		handlers:  make(map[string]CommandHandler),
		separator: " ",
	}
}

// AddCommandHandler adds a command handler
// for command
// e.g. for command "/help"
// handler should send help information to user
func (b *Bot) AddCommandHandler(handler CommandHandler, command string) {
	if _, ok := b.handlers[command]; ok {
		panic(fmt.Sprintf("Command handler with name %s already exists", command))
	}
	b.handlers[command] = handler
}

// AddCommandHandlerFunc adds a simple handler function
// for command without arguments
// e.g. for command "/help"
// handler should send help information to user
func (b *Bot) AddCommandHandlerFunc(function func(ctx context.Context, in *model.MessageIn) (out *model.MessageOut), command string) {
	_, ok := b.handlers[command]
	if !ok {
		b.handlers[command] = CommandHandlerFunc(function)
	}
}

func (b *Bot) Start() {
	if b.EnrichContext == nil {
		b.EnrichContext = GetContextFunc(func(_ *model.MessageIn) (context.Context, error) {
			return context.Background(), nil
		})
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updChan := b.Bot.GetUpdatesChan(updateConfig)

	for update := range updChan {
		switch {
		case update.Message != nil:
			b.HandleMessage(update.Message)
			break
		case update.CallbackQuery != nil:
			message := update.CallbackQuery.Message
			message.Text = update.CallbackQuery.Data
			message.From = update.CallbackQuery.From

			b.HandleMessage(message)
			break
		}
	}
}

func (b *Bot) SendResponse(in *model.MessageIn, out *model.MessageOut) {
	_, err := b.Bot.Send(out.ToTgMessage(in.Chat.ID))
	if err != nil {
		logrus.Infof("Ошибка при отправке сообщения: %v", err)
	}
	logrus.Infof("Пользователь %d написал %s и получил ответ %s",
		in.From.ID,
		in.Text,
		out.Text)
}

func (b *Bot) HandleMessage(in *tgbotapi.Message) {
	message := &model.MessageIn{
		Message: in,
	}
	var handler CommandHandler
	if strings.HasPrefix(message.Text, "/") {
		args := strings.Split(message.Text, b.separator)
		message.Command = args[0]
		if len(args) > 1 {
			message.Args = args[1:]
		}
		handler = b.handlers[message.Command]
		if handler == nil {
			b.SendResponse(message, &model.MessageOut{Text: fmt.Sprintf("Не знаю как обработать команду \"%s\"", message.Command)})
			return
		}
		b.chats[message.Chat.ID] = handler
	}

	ctx, err := b.GetContext(message)
	if err != nil {
		// TODO
		return
	}

	var messageOut *model.MessageOut

	if handler == nil {
		handler = b.chats[message.Chat.ID]
		if handler == nil {
			b.SendResponse(message, &model.MessageOut{Text: fmt.Sprintf("Не знаю что ответить на \"%s\"", message.Text)})
			return
		}
	}

	messageOut = handler.Handle(ctx, message)
	if messageOut != nil {
		b.SendResponse(message, messageOut)
	}
}

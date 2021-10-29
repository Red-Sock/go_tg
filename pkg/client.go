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
			b.handleMessage(update.Message)
			break
		case update.CallbackQuery != nil:
			message := update.CallbackQuery.Message
			message.Text = update.CallbackQuery.Data
			message.From = update.CallbackQuery.From

			b.handleMessage(message)
			break
		}
	}
}

func (b *Bot) SendMessage(t TgMessage, chatId int64) error {
	return t.Send(b.Bot, chatId)
}

func (b *Bot) sendResponse(in *model.MessageIn, out TgMessage) {
	err := out.Send(b.Bot, in.Chat.ID)
	if err != nil {
		logrus.Infof("Ошибка при отправке сообщения: %v", err)
	}
	logrus.Infof("Пользователь %d написал %s и получил ответ %v",
		in.From.ID,
		in.Text,
		out)
}

func (b *Bot) handleMessage(in *tgbotapi.Message) {
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
		handler = b.chats[message.Chat.ID]
		if handler != nil {
			handler.Dump(message.Chat.ID)
		}
		handler = b.handlers[message.Command]
		if handler == nil {
			b.sendResponse(message, &model.MessageOut{Text: fmt.Sprintf("Не знаю как обработать команду \"%s\"", message.Command)})
			return
		}
		b.chats[message.Chat.ID] = handler
	}

	ctx, err := b.GetContext(message)
	if err != nil {
		// TODO
		return
	}

	var messageOut TgMessage

	if handler == nil {
		handler = b.chats[message.Chat.ID]
		if handler == nil {
			b.sendResponse(message, &model.MessageOut{Text: fmt.Sprintf("Не знаю что ответить на \"%s\"", message.Text)})
			return
		}
	}

	messageOut = handler.Handle(ctx, message)
	if messageOut != nil {
		b.sendResponse(message, messageOut)
	}
}

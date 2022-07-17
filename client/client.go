package client

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/AlexSkilled/go_tg/handlers"
	"github.com/AlexSkilled/go_tg/interfaces"
	"github.com/AlexSkilled/go_tg/internal"
	"github.com/AlexSkilled/go_tg/model"
	"github.com/AlexSkilled/go_tg/model/menu"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// Bot - allows you to interact with telegram bot
// with some features
// tgbotapi.BotAPI - realisation of API calls to Telegram;
// chats - mapping of chat ids to their current handlers;
// handlers - mapping of name of handler to realisation;
// External context - can be used to pass information (such as user info) to handlers
// menuPattern - menu interaction(todo needs to be reworked)
type Bot struct {
	Bot *tgbotapi.BotAPI

	chats    map[int64]chatHandler
	handlers map[string]interfaces.CommandHandler

	interfaces.ExternalContext
	separator string

	menuPatterns    []interfaces.Menu
	locMenuPatterns []menu.LocalizedMenu
	menuHandler     *handlers.MenuHandler

	qm         *quitManager
	outMessage chan interfaces.Instruction
}

type quitManager struct {
	end chan struct{}
	wg  *sync.WaitGroup
}

// NewBot Bot constructor
func NewBot(token string) *Bot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &Bot{
		Bot:         bot,
		chats:       make(map[int64]chatHandler),
		handlers:    make(map[string]interfaces.CommandHandler),
		menuHandler: handlers.NewMenuHandler(),
		separator:   " ",
	}
}

// AddCommandHandler adds a command handler
// for command
// e.g. for command "/help"
// handler should send help information to user
func (b *Bot) AddCommandHandler(handler interfaces.CommandHandler, command string) {
	if _, ok := b.handlers[command]; ok {
		panic(fmt.Sprintf("Command handler with name %s already exists", command))
	}
	b.handlers[command] = handler
}

func (b *Bot) AddMenu(pattern interfaces.Menu) {
	b.menuHandler.AddSimpleMenu(pattern)
}

func (b *Bot) AddLocalizedMenu(locMenu menu.LocalizedMenu) {
	b.menuHandler.AddLocalizedMenu(locMenu)
}

func (b *Bot) Start() {
	// Context
	if b.ExternalContext == nil {
		b.ExternalContext = interfaces.GetContextFunc(func(_ *model.MessageIn) (context.Context, error) {
			return context.Background(), nil
		})
	}
	// menu handler at menu.MenuCall
	b.menuHandler.Retry = b.handleMessage
	b.handlers[menu.MenuCall] = b.menuHandler
	b.handlers[menu.Back] = b.menuHandler

	// Start
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updChan := b.Bot.GetUpdatesChan(updateConfig)

	quit := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(2)

	b.qm = &quitManager{
		quit,
		wg,
	}

	b.outMessage = make(chan interfaces.Instruction)
	internal.SetSender(b.outMessage)

	go b.handleInComing(updChan, b.qm)
	go b.handleOutgoing(b.qm)

}

func (b *Bot) Stop() {
	close(b.qm.end)
	b.qm.wg.Wait()
}

func (b *Bot) handleInComing(updChan tgbotapi.UpdatesChannel, qm *quitManager) {
	for {
		select {
		case update := <-updChan:
			switch {
			case update.Message != nil:

				b.handleMessage(&model.MessageIn{
					Message: update.Message,
				}, b.outMessage)
				break
			case update.CallbackQuery != nil:
				message := update.CallbackQuery.Message
				message.Text = update.CallbackQuery.Data
				message.From = update.CallbackQuery.From

				_, err := b.Bot.Request(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
				if err != nil {
					logrus.Error(err)
				}
				b.handleMessage(&model.MessageIn{
					Message: update.CallbackQuery.Message,
				}, b.outMessage)
				break
			}
		case <-qm.end:
			logrus.Println("Gracefully shutted down incoming handler")
			qm.wg.Done()
			return
		}

	}
}

func (b *Bot) handleOutgoing(qm *quitManager) {
	for {
		select {
		case inst := <-b.outMessage:
			// If outgoing message is menu - add menu
			switch t := inst.(type) {
			case interfaces.Menu:
				b.menuHandler.AttachMenu(inst.GetChatId(), t)
			case *model.RerenderMenu:
				inst = b.menuHandler.ReattachMenu(t)
				if inst == nil {
					continue
				}
				inst.SetMessageId(t.MessageId)
				inst.SetChatIdIfZero(t.ChatId)
			}

			sendMsg, err := b.Bot.Send(inst.GetMessage())
			if err != nil {
				logrus.Error(err)
			}

			inst.SetMessageId(int64(sendMsg.MessageID))

		case <-qm.end:
			logrus.Println("Gracefully shutted down outgoing handler")
			qm.wg.Done()
			return
		}
	}

}

func (b *Bot) handleMessage(message *model.MessageIn, outMessage chan<- interfaces.Instruction) {
	resp := &responser{c: outMessage, chatId: message.Chat.ID}

	ctx, err := b.GetContext(message)
	if err != nil {
		logrus.Error(err)
		b.outMessage <- model.NewMessageToChat(fmt.Sprintf("Couldn't GetContext, %v", err), message.Chat.ID)
		return
	}
	message.Ctx = ctx
	handler := b.chooseHandler(message)

	if handler == nil {
		msg := "Couldn't handle " + message.Command + " command"
		logrus.Error(msg)
		resp.Send(model.NewMessage(msg))
		return
	}

	handler.Handle(message, resp)
}

func (b *Bot) chooseHandler(message *model.MessageIn) (handler interfaces.CommandHandler) {
	if strings.HasPrefix(message.Text, "/") {
		// Extract command
		args := strings.Split(message.Text, b.separator)
		message.Command = args[0]
		if len(args) > 1 {
			message.Args = args[1:]
		}

		activeHandler, ok := b.chats[message.Chat.ID]
		if ok {
			// If message command and active handler's command are the same - use active in this chat handler
			if message.Command == activeHandler.command {
				return activeHandler.handler
			}
		}

		handler, ok = b.handlers[message.Command]
		if !ok && b.menuHandler.CanHandle(message) {
			// Dump session for this handler if message calls for a new command
			if b.menuHandler != activeHandler.handler && activeHandler.handler != nil {
				activeHandler.handler.Dump(message.Chat.ID)
			}
			handler = b.menuHandler
		}
		b.chats[message.Chat.ID] = chatHandler{command: message.Command, handler: handler}
		return handler
	}
	return b.chats[message.Chat.ID].handler
}

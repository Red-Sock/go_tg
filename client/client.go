package client

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/AlexSkilled/go_tg/handlers"
	"github.com/AlexSkilled/go_tg/interfaces"
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

	chats    map[int64]interfaces.CommandHandler
	handlers map[string]interfaces.CommandHandler

	interfaces.ExternalContext
	separator string

	menuPatterns    []interfaces.Menu
	locMenuPatterns []menu.LocalizedMenu
	menuHandler     *handlers.MenuHandler

	qm         *quitManager
	outMessage chan interfaces.Instruction
}

var instructionHandler <-chan interfaces.Instruction

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
		chats:       make(map[int64]interfaces.CommandHandler),
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
	b.handlers[menu.MenuCall] = b.menuHandler

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
				})
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
				})
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
				delete(b.chats, inst.GetChatId())
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

func (b *Bot) handleMessage(message *model.MessageIn) {

	var handler interfaces.CommandHandler

	resp := &Responser{c: b.outMessage, chatId: message.Chat.ID}

	if strings.HasPrefix(message.Text, "/") {
		args := strings.Split(message.Text, b.separator)
		message.Command = args[0]
		if len(args) > 1 {
			message.Args = args[1:]
		}
		handler = b.chats[message.Chat.ID]
		if handler != nil {
			handler.Dump(message.Chat.ID)
			delete(b.chats, message.Chat.ID)
		}
		handler = b.handlers[message.Command]

		b.chats[message.Chat.ID] = handler
	}

	ctx, err := b.GetContext(message)
	if err != nil {
		logrus.Error(err)
		b.outMessage <- &model.MessageOut{ChatId: message.Chat.ID, Text: fmt.Sprintf("Couldn't GetContext, %v", err)}
		return
	}

	message.Ctx = ctx

	if handler == nil {
		handler = b.chats[message.Chat.ID]
		if handler == nil {
			var ok bool
			handler, ok = b.handlers[menu.MenuCall]
			if !ok {
				logrus.Error("Couldn't handle", message.Args[0], "command")
				return
			}
		}
	}
	handler.Handle(message, resp)
}

// If there is NO handler for given command - tries to execute
// command with menu.MenuCall prefix
func (b *Bot) tryHandleAsMenuCall(in *model.MessageIn) {
	menuHandler, ok := b.handlers[menu.MenuCall]
	if !ok {
		return
	}
	menuHandler.Handle(in, &Responser{chatId: in.Chat.ID, c: b.outMessage})
}

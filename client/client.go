package client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/go_tg/handlers"
	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	menu2 "github.com/Red-Sock/go_tg/model/response/menu"
	"github.com/Red-Sock/go_tg/send"
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

	chats    map[int64]*chatHandler
	handlers map[string]interfaces.CommandHandler

	interfaces.ExternalContext
	separator string

	menuPatterns    []interfaces.Menu
	locMenuPatterns []menu2.LocalizedMenu
	menuHandler     *handlers.MenuHandler

	qm              *quitManager
	outMessage      chan interfaces.MessageOut
	responseTimeout time.Duration
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
		Bot:             bot,
		chats:           make(map[int64]*chatHandler),
		handlers:        make(map[string]interfaces.CommandHandler),
		menuHandler:     handlers.NewMenuHandler(),
		separator:       " ",
		responseTimeout: interfaces.UserResponseTimeout,
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

func (b *Bot) AddLocalizedMenu(locMenu menu2.LocalizedMenu) {
	b.menuHandler.AddLocalizedMenu(locMenu)
}

// SetResponseTimeout - sets timeout for user to response
// e.g. using interfaces.Chat's method GetInput will either wait for
// given @timeout or
// default timeout - interfaces.UserResponseTimeout or
// timeout provided via context
func (b *Bot) SetResponseTimeout(timeout time.Duration) {
	b.responseTimeout = timeout
}

func (b *Bot) Start() error {
	// Context
	if b.ExternalContext == nil {
		b.ExternalContext = interfaces.GetContextFunc(func(_ *model.MessageIn) (context.Context, error) {
			return context.Background(), nil
		})
	}
	// menu handler at menu.MenuCall
	b.menuHandler.Retry = b.handleMessage
	b.handlers[menu2.Back] = b.menuHandler

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

	b.outMessage = make(chan interfaces.MessageOut)
	send.SetSender(b.outMessage)

	defer func() {
		go b.handleInComing(updChan, b.qm)
		go b.handleOutgoing(b.qm)
	}()

	commands := make([]tgbotapi.BotCommand, 0, len(b.handlers)+len(b.menuPatterns))

	for command, handler := range b.handlers {
		err := validateCommand(command)
		if err != nil {
			return fmt.Errorf("error in command: %s. %w", command, err)
		}

		commands = append(commands, tgbotapi.BotCommand{
			Command:     command,
			Description: handler.GetDescription(),
		})
	}

	for _, item := range b.menuPatterns {
		command := item.GetCallCommand()

		err := validateCommand(command)
		if err != nil {
			return fmt.Errorf("error in command: %s. %w", command, err)
		}

		commands = append(commands, tgbotapi.BotCommand{
			Command:     command,
			Description: item.GetDescription(),
		})
	}

	rsp, err := b.Bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		return errors.Join(errors.New("error performing bot request to update commands"), err)
	}

	if !rsp.Ok {
		jsn, err := rsp.Result.MarshalJSON()
		if err != nil {
			return errors.Join(errors.New("error marshalling tg response"), err)
		}
		return errors.New(string(jsn))
	}

	return nil
}

func (b *Bot) Stop() {
	b.Bot.StopReceivingUpdates()

	close(b.qm.end)
	b.qm.wg.Wait()
}

func (b *Bot) Send(msg interfaces.MessageOut) {
	b.outMessage <- msg
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
			case *response.OpenMenu:
				inst = b.menuHandler.StartMenu(t.Msg, b.outMessage)
				if inst == nil {
					continue
				}
			}

			sendMsg, err := b.Bot.Send(inst.GetMessage())
			if err != nil {
				logrus.Error(err)
			}

			inst.ForceSetMessageId(int64(sendMsg.MessageID))

		case <-qm.end:
			logrus.Println("Gracefully shut down outgoing handler")
			qm.wg.Done()
			return
		}
	}

}

func (b *Bot) handleMessage(message *model.MessageIn, outMessage chan<- interfaces.MessageOut) {
	resp := &chat{
		chatId:  message.Chat.ID,
		cOut:    outMessage,
		timeout: b.responseTimeout,
	}

	ctx, err := b.GetContext(message)
	if err != nil {
		logrus.Error(err)
		b.outMessage <- response.NewMessageToChat(fmt.Sprintf("Couldn't GetContext, %v", err), message.Chat.ID)
		return
	}
	message.Ctx = ctx

	// if message starts with command separator - treat it as a start of new command
	if strings.HasPrefix(message.Text, "/") {
		cHandler := b.chooseHandler(message)
		if cHandler.handler != nil {
			resp.cIn = cHandler.msgCh
			if cHandler != nil {
				go cHandler.handler.Handle(message, resp)
			}
			return
		}
	} else {
		handler, ok := b.chats[message.Chat.ID]
		if ok {
			handler.msgCh <- message
			return
		}
	}

	msg := "Couldn't handle " + message.Command + " command"
	logrus.Error(msg)
	resp.SendMessage(response.NewMessage(msg))
}

func (b *Bot) chooseHandler(message *model.MessageIn) *chatHandler {
	args := strings.Split(message.Text, b.separator)
	message.Command = args[0]
	if len(args) > 1 {
		message.Args = args[1:]
	}

	handler, ok := b.handlers[message.Command]
	if !ok && b.menuHandler.CanHandle(message) {
		handler = b.menuHandler
	}

	activeHandler, ok := b.chats[message.Chat.ID]
	if ok {
		activeHandler.handler = handler
		close(activeHandler.msgCh)
		activeHandler.msgCh = make(chan *model.MessageIn)
	} else {
		b.chats[message.Chat.ID] = &chatHandler{
			handler: handler,
			msgCh:   make(chan *model.MessageIn),
		}
	}

	return b.chats[message.Chat.ID]
}

func validateCommand(command string) error {
	if len(command) < 2 {
		return errors.New("no name entered")
	}

	if command[0] != '/' {
		return errors.New("command has to start with \"/\" symbol")
	}

	availableRanges := [][]int32{
		{95, 95},
		{48, 57},
		{97, 122},
	}
	for _, s := range command[1:] {
		var hasHitRange = false
		for _, r := range availableRanges {
			if s >= r[0] && s <= r[1] {
				hasHitRange = true
				break
			}
		}
		if !hasHitRange {
			return errors.New("name contains \"" + string(s) + "\" symbol")
		}
	}

	return nil
}

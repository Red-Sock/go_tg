package go_tg

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/internal"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/send"
)

type TgApi interface {
	Start() error
	Stop()
	AddCommandHandler(handler interfaces.CommandHandler)
	SetDefaultCommandHandler(h interfaces.Handler)
}

// Bot - allows you to interact with telegram bot
// with some features
// tgbotapi.BotAPI - realisation of API calls to Telegram;
// chats - mapping of chat ids to their current handlers;
// handlers - mapping of name of handler to realisation;
// External context - can be used to pass information (such as user info) to handlers
// menuPattern - menu interaction(todo needs to be reworked)
type Bot struct {
	Bot *tgbotapi.BotAPI

	handlers       map[string]interfaces.CommandHandler
	defaultHandler interfaces.Handler

	interfaces.ExternalContext
	separator string

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
		handlers:        map[string]interfaces.CommandHandler{},
		separator:       " ",
		responseTimeout: interfaces.UserResponseTimeout,
		defaultHandler:  &internal.DefaultHandler{},
	}
}

// SetDefaultCommandHandler sets custom handler for unresolved messages
func (b *Bot) SetDefaultCommandHandler(h interfaces.Handler) {
	b.defaultHandler = h
}

// AddCommandHandler adds a command handler
// for command
// e.g. for command "/help"
// handler should send help information to user
func (b *Bot) AddCommandHandler(handler interfaces.CommandHandler) {
	command := handler.GetCommand()
	if _, ok := b.handlers[command]; ok {
		panic(fmt.Sprintf("Command handler with name %s already exists", command))
	}
	b.handlers[command] = handler
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
		b.ExternalContext = func(_ *model.MessageIn) context.Context {
			return context.Background()
		}
	}

	// HandlerMenu
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updChan := b.Bot.GetUpdatesChan(updateConfig)

	quit := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(1)

	b.qm = &quitManager{
		quit,
		wg,
	}

	send.SetSender(b.handleOutgoing)

	defer func() {
		go b.handleInComing(updChan, b.qm)
	}()

	commands := make([]tgbotapi.BotCommand, 0, len(b.handlers))

	for command, handler := range b.handlers {
		err := validateCommand(command)
		if err != nil {
			return fmt.Errorf("error in command: %s. %w", command, err)
		}

		d, ok := handler.(interfaces.Description)
		if ok {
			commands = append(commands, tgbotapi.BotCommand{
				Command:     command,
				Description: d.GetDescription(),
			})
		}
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
	b.handleOutgoing(msg)
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

			case update.CallbackQuery != nil:
				message := update.CallbackQuery.Message
				message.Text = update.CallbackQuery.Data
				message.From = update.CallbackQuery.From

				_, err := b.Bot.Request(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
				if err != nil {
					logrus.Errorf("error responsing to callback %s", err)
				}
				b.handleMessage(&model.MessageIn{
					Message:    update.CallbackQuery.Message,
					IsCallback: true,
				})

			}
		case <-qm.end:
			logrus.Println("Gracefully shutted down incoming handler")
			qm.wg.Done()
			return
		}
	}
}

func (b *Bot) handleOutgoing(out interfaces.MessageOut) error {
	sendMsg, err := b.Bot.Send(out.GetMessage())
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "json: cannot unmarshal bool into Go value of type tgbotapi.Message") {
			return nil
		}
		if strings.Contains(errMsg, "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message") {
			return nil
		}
		return err
	}

	out.ForceSetMessageId(int64(sendMsg.MessageID))

	return nil
}

func (b *Bot) handleMessage(message *model.MessageIn) {
	resp := &internal.Chat{
		ChatId:  message.Chat.ID,
		COut:    b.handleOutgoing,
		Timeout: b.responseTimeout,
	}
	message.Args = strings.Split(message.Text, " ")

	if len(message.Args) != 0 && len(message.Args[0]) != 0 {
		if message.Args[0][0] == '/' {
			message.Command = message.Args[0]
			message.Args = message.Args[1:]
		}
	}

	message.Ctx = b.ExternalContext(message)

	var handler interfaces.Handler

	handler, ok := b.handlers[message.Command]
	if !ok {
		handler = b.defaultHandler
	}

	logrus.Infof("%s with args %v", message.Command, message.Args)

	handler.Handle(message, resp)
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

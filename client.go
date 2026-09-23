package go_tg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

	MustAddCommandHandler(handler interfaces.CommandHandler)
	AddCommandHandler(handler interfaces.CommandHandler) error

	SetDefaultCommandHandler(h interfaces.Handler)
	SetInlineQueryHandler(h interfaces.InlineQueryHandler)
	SetChosenInlineResultHandler(h interfaces.ChosenInlineResultHandler)

	Send(msg interfaces.MessageOut) error
	SendAndReturn(msg interfaces.MessageOut) (tgbotapi.Message, error)
	EditInlineMessageAudio(inlineMessageId, fileId string) error
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

	handlers                  map[string]interfaces.CommandHandler
	defaultHandler            interfaces.Handler
	inlineQueryHandler        interfaces.InlineQueryHandler
	chosenInlineResultHandler interfaces.ChosenInlineResultHandler

	interfaces.ExternalContext
	separator string

	qm              *quitManager
	outMessage      chan interfaces.MessageOut
	responseTimeout time.Duration

	middlewares     []interfaces.Middleware
	logger          logrus.FieldLogger
	onlyDirectCalls bool
	httpClient      tgbotapi.HTTPClient
}

type quitManager struct {
	end chan struct{}
	wg  *sync.WaitGroup
}

// NewBot Bot constructor
func NewBot(token string, opts ...opt) (*Bot, error) {
	botInstance := &Bot{
		handlers:        map[string]interfaces.CommandHandler{},
		separator:       " ",
		responseTimeout: interfaces.UserResponseTimeout,
		logger:          logrus.New(),
	}

	for _, o := range opts {
		o(botInstance)
	}

	if botInstance.httpClient == nil {
		botInstance.httpClient = &loggingHTTPClient{
			inner:  http.DefaultClient,
			logger: botInstance.logger,
		}
	}

	botApi, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, botInstance.httpClient)
	if err != nil {
		return nil, fmt.Errorf("error creating tg bot connection, %w", err)
	}

	botInstance.Bot = botApi

	botInstance.defaultHandler = &internal.DefaultHandler{
		Logger: botInstance.logger,
	}

	return botInstance, nil
}

// Use registers one or more middlewares. They run before every handler,
// in registration order (first registered = outermost = runs first).
func (b *Bot) Use(mw ...interfaces.Middleware) {
	b.middlewares = append(b.middlewares, mw...)
}

// SetDefaultCommandHandler sets custom handler for unresolved messages
func (b *Bot) SetDefaultCommandHandler(h interfaces.Handler) {
	b.defaultHandler = h
}

// SetInlineQueryHandler sets the handler for inline queries - triggered when
// a user types "@botname <query>" in any chat. A bot without one set simply
// ignores inline queries (Telegram shows no results, no error).
func (b *Bot) SetInlineQueryHandler(h interfaces.InlineQueryHandler) {
	b.inlineQueryHandler = h
}

// SetChosenInlineResultHandler sets the handler notified when a user taps one
// of the results a previous inline query was answered with. Requires
// inline feedback enabled for the bot via BotFather (/setinlinefeedback) -
// Telegram otherwise never sends chosen_inline_result updates.
func (b *Bot) SetChosenInlineResultHandler(h interfaces.ChosenInlineResultHandler) {
	b.chosenInlineResultHandler = h
}

// AddCommandHandler adds a command handler
// for command
// e.g. for command "/help"
// handler should send help information to user
func (b *Bot) AddCommandHandler(handler interfaces.CommandHandler) error {
	command := handler.GetCommand()
	if _, ok := b.handlers[command]; ok {
		return fmt.Errorf("command handler with name %s already exists", command)
	}

	b.handlers[command] = handler

	return nil
}
func (b *Bot) MustAddCommandHandler(handler interfaces.CommandHandler) {
	err := b.AddCommandHandler(handler)
	if err != nil {
		panic(err)
	}
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

	commandsWithDescription := make([]tgbotapi.BotCommand, 0, len(b.handlers))

	for command, handler := range b.handlers {
		err := validateCommand(command)
		if err != nil {
			return fmt.Errorf("error in command: %s. %w", command, err)
		}

		d, ok := handler.(interfaces.Description)
		if ok {
			commandsWithDescription = append(commandsWithDescription, tgbotapi.BotCommand{
				Command:     command,
				Description: d.GetDescription(),
			})
		}
	}

	rsp, err := b.Bot.Request(tgbotapi.NewSetMyCommands(commandsWithDescription...))
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

	b.handleInComing(updChan, b.qm)
	return nil
}

func (b *Bot) Stop() {
	b.Bot.StopReceivingUpdates()

	close(b.qm.end)
	b.qm.wg.Wait()
}

func (b *Bot) Send(msg interfaces.MessageOut) error {
	err := b.handleOutgoing(msg)
	if err != nil {
		b.logger.WithError(err).Error("error handling outgoing message")
		return err
	}

	return nil
}

// SendAndReturn behaves like Send but also hands back Telegram's response
// message, e.g. to read the file_id Telegram assigned a freshly uploaded
// audio file for later reuse via an InlineQueryResultCached* result.
func (b *Bot) SendAndReturn(msg interfaces.MessageOut) (tgbotapi.Message, error) {
	sendMsg, err := b.sendAndReturn(msg)
	if err != nil {
		b.logger.WithError(err).Error("error handling outgoing message")
		return tgbotapi.Message{}, err
	}

	return sendMsg, nil
}

// EditInlineMessageAudio replaces an inline message's content with the audio
// identified by fileId - used to redeliver an already-known Telegram file_id
// into whichever chat the inline message lives in, once it's ready.
func (b *Bot) EditInlineMessageAudio(inlineMessageId, fileId string) error {
	edit := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{InlineMessageID: inlineMessageId},
		Media:    tgbotapi.NewInputMediaAudio(tgbotapi.FileID(fileId)),
	}

	_, err := b.Bot.Request(edit)
	if err != nil {
		b.logger.WithError(err).Error("error editing inline message audio")
		return err
	}

	return nil
}

func (b *Bot) handleInComing(updChan tgbotapi.UpdatesChannel, qm *quitManager) {
	for {
		select {
		case update := <-updChan:
			switch {
			case update.Message != nil:
				msg := &model.MessageIn{
					Message: update.Message,
				}
				b.handleMessage(msg)

			case update.CallbackQuery != nil:
				message := update.CallbackQuery.Message
				message.Text = update.CallbackQuery.Data
				message.From = update.CallbackQuery.From

				_, err := b.Bot.Request(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
				if err != nil {
					b.logger.Errorf("error responsing to callback %s", err)
				}

				msg := &model.MessageIn{
					Message:    update.CallbackQuery.Message,
					IsCallback: true,
				}

				b.handleMessage(msg)

			case update.InlineQuery != nil:
				b.handleInlineQuery(update.InlineQuery)

			case update.ChosenInlineResult != nil:
				b.handleChosenInlineResult(update.ChosenInlineResult)
			}
		case <-qm.end:
			b.logger.Println("Gracefully shutted down incoming handler")
			qm.wg.Done()
			return
		}
	}
}

func (b *Bot) handleOutgoing(out interfaces.MessageOut) error {
	_, err := b.sendAndReturn(out)
	return err
}

// sendAndReturn sends out and, unlike handleOutgoing, also hands back
// Telegram's response message - needed by SendAndReturn to read back fields
// Telegram assigns on send (e.g. Audio.FileID for a freshly uploaded file).
func (b *Bot) sendAndReturn(out interfaces.MessageOut) (tgbotapi.Message, error) {
	sendMsg, err := b.Bot.Send(out.GetMessage())
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "json: cannot unmarshal bool into Go value of type tgbotapi.Message") {
			return tgbotapi.Message{}, nil
		}
		if strings.Contains(errMsg, "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message") {
			return tgbotapi.Message{}, nil
		}

		return tgbotapi.Message{}, err
	}

	out.ForceSetMessageId(int64(sendMsg.MessageID))

	return sendMsg, nil
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

			if strings.HasSuffix(message.Command, "@"+b.Bot.Self.UserName) {
				message.Command = message.Command[:strings.LastIndex(message.Command, "@")]
			} else if b.onlyDirectCalls {
				return
			}
		}
	}

	message.Ctx = b.ExternalContext(message)

	var handler interfaces.Handler

	handler, ok := b.handlers[message.Command]
	if !ok {
		handler = b.defaultHandler
	}

	for i := len(b.middlewares) - 1; i >= 0; i-- {
		handler = b.middlewares[i](handler)
	}

	err := handler.Handle(message, resp)
	if err != nil {
		b.logger.Errorf("%s with args %v error: %v", message.Command, message.Args, err)
	} else {
		b.logger.Infof("%s with args %v", message.Command, message.Args)
	}
}

func (b *Bot) handleInlineQuery(q *tgbotapi.InlineQuery) {
	if b.inlineQueryHandler == nil {
		return
	}

	in := &model.InlineQueryIn{InlineQuery: q}

	results, err := b.inlineQueryHandler.Handle(in)
	if err != nil {
		b.logger.Errorf("inline query %q error: %v", q.Query, err)
		return
	}

	answer := tgbotapi.InlineConfig{
		InlineQueryID: q.ID,
		Results:       results,
	}

	_, err = b.Bot.Request(answer)
	if err != nil {
		b.logger.Errorf("error answering inline query %q: %v", q.Query, err)
	}
}

func (b *Bot) handleChosenInlineResult(r *tgbotapi.ChosenInlineResult) {
	if b.chosenInlineResultHandler == nil {
		return
	}

	in := &model.ChosenInlineResultIn{ChosenInlineResult: r}

	err := b.chosenInlineResultHandler.HandleChosen(in)
	if err != nil {
		b.logger.Errorf("chosen inline result %q error: %v", r.ResultID, err)
	}
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

package interfaces

import (
	"github.com/Red-Sock/go_tg/model"
)

// ChosenInlineResultHandler handles a Telegram chosen inline result (see
// model.ChosenInlineResultIn) - notified after a user taps one of the
// results a previous InlineQueryHandler.Handle call answered with.
type ChosenInlineResultHandler interface {
	HandleChosen(in *model.ChosenInlineResultIn) error
}

// ChosenInlineResultHandlerFunc is a function adapter for
// ChosenInlineResultHandler.
type ChosenInlineResultHandlerFunc func(in *model.ChosenInlineResultIn) error

func (f ChosenInlineResultHandlerFunc) HandleChosen(in *model.ChosenInlineResultIn) error {
	return f(in)
}

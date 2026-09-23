package interfaces

import (
	"github.com/Red-Sock/go_tg/model"
)

// InlineQueryHandler answers a Telegram inline query (see model.InlineQueryIn).
// The []interface{} return mirrors tgbotapi.InlineConfig.Results: Telegram's
// Bot API has no common Go type for its ~15 inline result kinds
// (InlineQueryResultArticle, InlineQueryResultAudio, ...), so tgbotapi itself
// takes them as interface{}.
type InlineQueryHandler interface {
	Handle(in *model.InlineQueryIn) ([]interface{}, error)
}

// InlineQueryHandlerFunc is a function adapter for InlineQueryHandler.
type InlineQueryHandlerFunc func(in *model.InlineQueryIn) ([]interface{}, error)

func (f InlineQueryHandlerFunc) Handle(in *model.InlineQueryIn) ([]interface{}, error) {
	return f(in)
}

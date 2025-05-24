package go_tg

import (
	"github.com/sirupsen/logrus"
)

type opt func(bot *Bot)

func WithLogger(logger logrus.FieldLogger) opt {
	return func(bot *Bot) {
		bot.logger = logger
	}
}

// WithOnlyDirectCalls - set value for "Direct calls" flag
// When "Direct calls" flag is set in group chats bot will
// be trigger only when command is called with bot tag
//
//	e.g.
//		with onlyDirectCalls = true
//		command /version - won't be handled by bot
//		but command /verstion@bot_name_here_bot - will trigger handler
func WithOnlyDirectCalls(v bool) opt {
	return func(bot *Bot) {
		bot.onlyDirectCalls = v
	}
}

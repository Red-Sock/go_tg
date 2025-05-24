package go_tg

import (
	"github.com/sirupsen/logrus"
)

type opt func(bot *Bot)

func SetLogger(logger logrus.FieldLogger) opt {
	return func(bot *Bot) {
		bot.logger = logger
	}
}

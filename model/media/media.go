package media

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Media interface {
	AsSingleTgMedia(chatId int64) tgbotapi.Chattable
	AsInputMedia() any
}

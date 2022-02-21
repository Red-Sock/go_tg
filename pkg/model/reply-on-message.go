package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Reply struct {
	TgCallbackQueryId string
	Text              string
	ShowAlert         bool
	ShowTime          int
}

func (r *Reply) Send(api *tgbotapi.BotAPI, _ int64) error {
	replyMessage := tgbotapi.NewCallback(r.TgCallbackQueryId, r.Text)

	replyMessage.ShowAlert = r.ShowAlert

	if r.ShowTime == 0 {
		replyMessage.CacheTime = 5
	} else {
		replyMessage.CacheTime = r.ShowTime
	}

	_, err := api.Send(replyMessage)
	return err
}

package media

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Video struct {
	Caption string
	FileID  string
}

func (m Video) AsSingleTgMedia(chatId int64) tgbotapi.Chattable {
	vid := tgbotapi.NewVideo(chatId, tgbotapi.FileID(m.FileID))
	vid.Caption = m.Caption

	return vid
}

func (m Video) AsInputMedia() any {
	return tgbotapi.NewInputMediaVideo(tgbotapi.FileID(m.FileID))
}

package response

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/AlexSkilled/go_tg/model"
)

type OpenMenu struct {
	Msg *model.MessageIn
}

func NewOpenMenu(menuCommand string, originalMsg *model.MessageIn) *OpenMenu {
	originalMsg.Command = menuCommand
	return &OpenMenu{
		Msg: originalMsg,
	}
}

func (u *OpenMenu) GetMessage() tgbotapi.Chattable {
	return nil
}

func (u *OpenMenu) GetChatId() int64 {
	return u.Msg.Chat.ID
}

func (u *OpenMenu) SetChatIdIfZero(c int64) {
	if u.Msg.Chat.ID == 0 {
		u.Msg.Chat.ID = c
	}
}

func (u *OpenMenu) ForceSetMessageId(id int64) {
	u.Msg.MessageID = int(id)
}

func (u *OpenMenu) GetMessageId() int64 {
	return int64(u.Msg.MessageID)
}

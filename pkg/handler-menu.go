package pkg

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlexSkilled/go_tg/pkg/model"
	"github.com/sirupsen/logrus"
)

type MenuHandler struct {
	patterns   map[string]model.Menu // command to pattern
	chatToMenu map[int64]model.Menu  // chat id to menu
}

func newMenuHandler(bot *Bot) *MenuHandler {
	mh := &MenuHandler{
		patterns:   map[string]model.Menu{},
		chatToMenu: map[int64]model.Menu{},
	}

	for _, item := range bot.menuPatterns {
		menu := model.NewMenu(item)
		mh.patterns[menu.GetCallCommand()] = menu
	}

	return mh
}

func (m *MenuHandler) Handle(ctx context.Context, in *model.MessageIn) (out TgMessage) {
	menu, ok := m.chatToMenu[in.Chat.ID]

	switch in.Command {
	case model.MenuCall:
		return m.handleMenuCall(ctx, in)
	case model.ChangePage:
		if ok {
			var page int
			var err error

			if len(in.Args) > 0 {
				page, err = strconv.Atoi(in.Args[0])
				if err != nil {
					page = 0
				}
			}

			ik := menu.GetPage(ctx, page)
			return model.EditMessageReply(&ik, in.MessageID)
		}
	case model.OpenMenu:
		if len(in.Args) == 0 {
			logrus.Error(fmt.Sprintf("Error when tried to open new menu. No menu name specified"))
			return nil
		}
		in.Command = in.Args[0]
		return m.startMenu(ctx, in)
	default:
		return m.startMenu(ctx, in)
	}

	m.chatToMenu[in.Chat.ID] = menu
	return
}

func (m *MenuHandler) startMenu(ctx context.Context, in *model.MessageIn) TgMessage {
	if pattern, ok := m.patterns[in.Command]; ok {

		ik := pattern.GetPage(ctx, 0)
		return &model.MessageOut{
			Text:          pattern.GetName(ctx),
			InlineButtons: &ik,
		}
	}
	return nil
}

func (m *MenuHandler) handleMenuCall(ctx context.Context, in *model.MessageIn) TgMessage {
	if pattern, ok := m.patterns[in.Args[0]]; ok {
		var page int
		if len(in.Args) > 1 {
			page, _ = strconv.Atoi(in.Args[1])
		}

		ik := pattern.GetPage(ctx, page)
		name := pattern.GetName(ctx)
		return &model.MessageEdit{
			MessageId:     in.MessageID,
			Text:          &name,
			ReplyKeyboard: &ik,
		}
	}
	return nil
}

func (m *MenuHandler) Dump(id int64) {
	delete(m.chatToMenu, id)
}

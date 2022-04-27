package pkg

import (
	"strconv"

	"github.com/AlexSkilled/go_tg/pkg/model"
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

func (m *MenuHandler) Handle(in *model.MessageIn, out Sender) {
	switch in.Command {
	case model.MenuCall:
		out.Send(m.handleMenuCall(in))
	case model.ChangePage:
		if menu, ok := m.chatToMenu[in.Chat.ID]; ok {
			var page int
			var err error

			if len(in.Args) > 0 {
				page, err = strconv.Atoi(in.Args[0])
				if err != nil {
					page = 0
				}
			}

			ik := menu.GetPage(in.Ctx, page)
			model.EditMessageReply(&ik, in.MessageID)
		}
		out.Send(&model.MessageOut{
			Text: "No active menu for changing page",
		})
	default:
		out.Send(m.startMenu(in))
	}
}

func (m *MenuHandler) startMenu(in *model.MessageIn) Instruction {
	if len(in.Args) == 0 {
		return &model.MessageOut{
			Text: "Expected name of menu, but got nothing: " + model.MenuCall + " " + model.OpenMenu + " *Empty_Menu_Name*",
		}
	}
	if pattern, ok := m.patterns[in.Args[0]]; ok {
		ik := pattern.GetPage(in.Ctx, 0)

		m.chatToMenu[in.Chat.ID] = pattern

		return &model.MessageOut{
			Text:          pattern.GetName(in.Ctx),
			InlineButtons: &ik,
		}
	}
	return &model.MessageOut{
		Text: "Couldn't find registered menu with name: " + in.Args[0],
	}
}

func (m *MenuHandler) handleMenuCall(in *model.MessageIn) Instruction {
	if len(in.Args) == 0 {
		return nil
	}
	in.Command = in.Args[0]
	in.Args = in.Args[0:]
	if in.Command == model.OpenMenu {
		return m.startMenu(in)

	}

	if pattern, ok := m.patterns[in.Args[0]]; ok {
		var page int
		if len(in.Args) > 1 {
			page, _ = strconv.Atoi(in.Args[1])
		}

		ik := pattern.GetPage(in.Ctx, page)
		name := pattern.GetName(in.Ctx)
		pattern.SetPreviousMenu(m.chatToMenu[in.Chat.ID])
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

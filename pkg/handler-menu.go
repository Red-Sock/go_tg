package pkg

import (
	"fmt"
	"strconv"

	"github.com/AlexSkilled/go_tg/pkg/model"
)

type MenuHandler struct {
	simplePatterns    map[string]model.Menu            // command to menu
	localizedPatterns map[string]map[string]model.Menu // locale -> command -> menu
	chatToMenu        map[int64]model.Menu             // chat id to menu
}

func newMenuHandler(simpleMenus []model.Menu) *MenuHandler {
	mh := &MenuHandler{
		simplePatterns: map[string]model.Menu{},
		chatToMenu:     map[int64]model.Menu{},
	}

	for _, item := range simpleMenus {
		_, ok := mh.simplePatterns[item.GetCallCommand()]
		if ok {
			panic(fmt.Sprintf("Same menu for command %s", item.GetCallCommand()))
		}
		mh.simplePatterns[item.GetCallCommand()] = item
	}

	return mh
}

func (m *MenuHandler) AddLocalizedMenus(localizedMenu []model.LocalizedMenu) {
	m.localizedPatterns = map[string]map[string]model.Menu{}
	for _, item := range localizedMenu {
		for lang, menu := range item.LangToMenu {
			if in, ok := m.localizedPatterns[lang]; ok {
				in[menu.GetCallCommand()] = menu
			} else {
				m.localizedPatterns[lang] = map[string]model.Menu{
					menu.GetCallCommand(): menu,
				}
			}
		}
	}
}

func (m *MenuHandler) Handle(in *model.MessageIn, out Sender) {
	var menuStorage map[string]model.Menu
	var ok bool

	locale, ok := in.Ctx.Value(model.LocaleContextKey).(string)
	if ok {
		menuStorage, ok = m.localizedPatterns[locale]
	}
	if !ok {
		menuStorage = m.simplePatterns
	}

	switch in.Command {
	case model.MenuCall:
		out.Send(m.handleMenuCall(in, menuStorage))
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

			ik := menu.GetPage(page)
			model.EditMessageReply(ik, in.MessageID)
		}
		out.Send(&model.MessageOut{
			Text: "No active menu for changing page",
		})
	default:
		out.Send(m.startMenu(in, menuStorage))
	}
}

func (m *MenuHandler) startMenu(in *model.MessageIn, mStorage map[string]model.Menu) Instruction {
	if len(in.Args) == 0 {
		return &model.MessageOut{
			Text: "Expected name of menu, but got nothing: " + model.MenuCall + " " + model.OpenMenu + " *Empty_Menu_Name*",
		}
	}
	if pattern, ok := mStorage[in.Args[0]]; ok {
		m.chatToMenu[in.Chat.ID] = pattern

		return &model.MessageOut{
			Text:          pattern.GetName(),
			InlineButtons: pattern.GetPage(),
		}
	}
	return &model.MessageOut{
		Text: "Couldn't find registered menu with name: " + in.Args[0],
	}
}

func (m *MenuHandler) handleMenuCall(in *model.MessageIn, mStorage map[string]model.Menu) Instruction {
	if len(in.Args) == 0 {
		return nil
	}
	in.Command = in.Args[0]
	in.Args = in.Args[0:]
	if in.Command == model.OpenMenu {
		return m.startMenu(in, mStorage)

	}

	if pattern, ok := mStorage[in.Args[0]]; ok {
		var page int
		if len(in.Args) > 1 {
			page, _ = strconv.Atoi(in.Args[1])
		}

		name := pattern.GetName()
		pattern.SetPreviousMenu(m.chatToMenu[in.Chat.ID])
		return &model.MessageEdit{
			MessageId:     in.MessageID,
			Text:          &name,
			ReplyKeyboard: pattern.GetPage(page),
		}
	}
	return nil
}

func (m *MenuHandler) Dump(id int64) {
	delete(m.chatToMenu, id)
}

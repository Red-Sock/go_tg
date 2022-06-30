package handlers

import (
	"fmt"

	"github.com/AlexSkilled/go_tg/interfaces"
	"github.com/AlexSkilled/go_tg/model"
	"github.com/AlexSkilled/go_tg/model/menu"
)

type MenuHandler struct {
	simplePatterns    map[string]interfaces.Menu            // command to menu
	localizedPatterns map[string]map[string]interfaces.Menu // locale -> command -> menu
	chatToMenu        map[int64]interfaces.Menu             // chat id to menu
}

func NewMenuHandler() *MenuHandler {
	mh := MenuHandler{
		simplePatterns:    map[string]interfaces.Menu{},
		chatToMenu:        map[int64]interfaces.Menu{},
		localizedPatterns: map[string]map[string]interfaces.Menu{},
	}

	return &mh
}

func (m *MenuHandler) AddSimpleMenu(item interfaces.Menu) {
	_, ok := m.simplePatterns[item.GetCallCommand()]
	if ok {
		panic(fmt.Sprintf("Same menu for command %s", item.GetCallCommand()))
	}
	m.simplePatterns[item.GetCallCommand()] = item
}

func (m *MenuHandler) AddLocalizedMenu(item menu.LocalizedMenu) {
	for lang, menuItem := range item.LangToMenu {

		if in, ok := m.localizedPatterns[lang]; ok {
			in[menuItem.GetCallCommand()] = menuItem
		} else {
			m.localizedPatterns[lang] = map[string]interfaces.Menu{
				menuItem.GetCallCommand(): menuItem,
			}
		}
	}
}

func (m *MenuHandler) Handle(in *model.MessageIn, out interfaces.Sender) {
	var menuStorage map[string]interfaces.Menu
	var ok bool

	locale, ok := in.Ctx.Value(model.LocaleContextKey).(string)
	if ok {
		menuStorage, ok = m.localizedPatterns[locale]
	}
	if !ok {
		menuStorage = m.simplePatterns
	}

	switch in.Command {
	case menu.MenuCall:
		out.Send(m.handleMenuCall(in, menuStorage))
	case menu.Back:
		if currentMenu, ok := m.chatToMenu[in.Chat.ID]; ok {
			prev := currentMenu.GetPreviousMenu()
			if prev == nil {
				out.Send(&model.MessageOut{Text: "Nowhere to return!"})
				return
			}
			out.Send(prev)
		}
	//case menu.ChangePage:
	//	if menu, ok := m.chatToMenu[in.Chat.ID]; ok {
	//		var page int
	//		var err error
	//
	//		if len(in.Args) > 0 {
	//			page, err = strconv.Atoi(in.Args[0])
	//			if err != nil {
	//				page = 0
	//			}
	//		}
	//
	//		ik := menu.GetPage(page)
	//		model.EditMessageReply(ik, in.MessageID)
	//	}
	//	out.Send(&model.MessageOut{
	//		Text: "No active menu for changing page",
	//	})
	default:
		out.Send(m.startMenu(in, menuStorage))
	}
}

func (m *MenuHandler) startMenu(in *model.MessageIn, mStorage map[string]interfaces.Menu) interfaces.Instruction {
	if pattern, ok := mStorage[in.Text]; ok {

		if msg, ok := m.chatToMenu[in.Chat.ID]; ok {
			pattern = pattern.GetCopy()
			pattern.SetMessageId(msg.GetMessageId())
			pattern.SetPreviousMenu(msg)
			m.chatToMenu[in.Chat.ID] = pattern
		} else {
			pattern.SetMessageId(int64(in.MessageID))
		}
		pattern.SetChatIdIfZero(in.Chat.ID)

		return pattern
	}
	return &model.MessageOut{
		Text: "Couldn't find registered menu with name: " + in.Text,
	}
}

func (m *MenuHandler) handleMenuCall(in *model.MessageIn, mStorage map[string]interfaces.Menu) interfaces.Instruction {
	if len(in.Args) == 0 {
		return nil
	}
	in.Command = in.Args[0]
	in.Args = in.Args[0:]
	if in.Command == "model.OpenMenu" {
		return m.startMenu(in, mStorage)

	}

	if pattern, ok := mStorage[in.Args[0]]; ok {
		//var page int
		if len(in.Args) > 1 {
			//page, _ = strconv.Atoi(in.Args[1])
		}

		pattern.SetPreviousMenu(m.chatToMenu[in.Chat.ID])
		pattern.SetMessageId(int64(in.MessageID))

		return pattern
	}
	return nil
}

func (m *MenuHandler) Dump(id int64) {
	delete(m.chatToMenu, id)
}

func (m *MenuHandler) AttachMenu(chatId int64, menu interfaces.Menu) {
	m.chatToMenu[chatId] = menu
}

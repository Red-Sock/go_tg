package handlers

import (
	"context"
	"fmt"

	"github.com/AlexSkilled/go_tg/interfaces"
	"github.com/AlexSkilled/go_tg/model"
	"github.com/AlexSkilled/go_tg/model/menu"
)

const (
	NoMenuError = "Couldn't find registered menu with name: "
	defaultLang = "no_language_specified"
)

type MenuHandler struct {
	localizedPatterns map[string]map[string]interfaces.Menu // locale -> command -> menu
	chatToMenu        map[int64]interfaces.Menu             // chat id to menu

	customMenus map[int64]map[string]interfaces.Menu // userID to command to menu

	// Retry is a function that retries to handler message.
	// Exists for one purpose only - in case if (all AND conditions)
	// 1. backend has returned a custom created menu
	// 2. user has felt into sub menu
	// 3. user wants to return from sub menu to custom created menu
	// framework retries to call backend same command specified for parent menu
	Retry func(message *model.MessageIn, c chan<- interfaces.Instruction)
}

func NewMenuHandler() *MenuHandler {
	mh := MenuHandler{
		chatToMenu:        map[int64]interfaces.Menu{},
		localizedPatterns: map[string]map[string]interfaces.Menu{},
	}

	return &mh
}

func (m *MenuHandler) AddSimpleMenu(item interfaces.Menu) {
	_, ok := m.localizedPatterns[defaultLang][item.GetCallCommand()]
	if ok {
		panic(fmt.Sprintf("Same menu for command %s", item.GetCallCommand()))
	}
	m.localizedPatterns[defaultLang][item.GetCallCommand()] = item
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
	menuStorage := m.getMenuStorage(in.Ctx)

	switch in.Command {
	case menu.MenuCall:
		out.Send(m.handleMenuCall(in, menuStorage))
	case menu.Back:
		if currentMenu, ok := m.chatToMenu[in.Chat.ID]; ok {
			prev := currentMenu.GetPreviousMenu()
			if prev == nil {
				out.Send(model.NewMessage("Nowhere to return!"))
				return
			}
			ms := m.getMenuStorage(in.Ctx)

			newPrev := ms[prev.GetCallCommand()]
			if newPrev == nil {
				in.Text = prev.GetCallCommand()
				c := make(chan interfaces.Instruction)
				go m.Retry(in, c)
				inst := <-c

				var ok bool
				if newPrev, ok = inst.(interfaces.Menu); !ok {
					return
				}
			}
			newPrev.SetMessageId(prev.GetMessageId())
			newPrev.SetChatIdIfZero(prev.GetChatId())
			newPrev.SetPreviousMenu(prev.GetPreviousMenu())
			out.Send(newPrev)
		}
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
	return model.NewMessage(NoMenuError + in.Text)
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

func (m *MenuHandler) getMenuStorage(ctx context.Context) (menuStorage map[string]interfaces.Menu) {
	var ok bool

	locale, ok := ctx.Value(model.LocaleContextKey).(string)
	if !ok {
		locale = defaultLang
	}
	return m.localizedPatterns[locale]
}

func (m *MenuHandler) Dump(id int64) {
	delete(m.chatToMenu, id)
}

func (m *MenuHandler) AttachMenu(chatId int64, menu interfaces.Menu) {
	m.chatToMenu[chatId] = menu
}

func (m *MenuHandler) ReattachMenu(req *model.RerenderMenu) interfaces.Menu {
	currentMenu := m.chatToMenu[req.ChatId]
	if currentMenu == nil {
		return nil
	}
	newMenu := m.getMenuStorage(req.Ctx)[currentMenu.GetCallCommand()]
	newMenu.SetPreviousMenu(currentMenu.GetPreviousMenu())
	newMenu.SetMessageId(req.MessageId)
	m.chatToMenu[req.ChatId] = newMenu
	return newMenu
}

func (m *MenuHandler) CanHandle(in *model.MessageIn) bool {
	// Can handler if there is a chat with menu
	if _, ok := m.chatToMenu[in.Chat.ID]; ok {
		return true
	}

	return false
}

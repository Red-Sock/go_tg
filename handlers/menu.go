package handlers

import (
	"context"
	"fmt"

	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
	menu2 "github.com/Red-Sock/go_tg/model/response/menu"
)

const (
	NoMenuError = "Couldn't find registered menu with name: "
	defaultLang = "no_language_specified"
)

type MenuHandler struct {
	localizedPatterns map[string]map[string]interfaces.Menu // locale -> command -> menu
	chatToMenu        map[int64]interfaces.Menu             // chat id to menu

	// Retry is a function that retries to handler message.
	// Exists for one purpose only - in case if (all AND conditions)
	// 1. backend has returned a custom created menu
	// 2. user has felt into sub menu
	// 3. user wants to return from sub menu to custom created menu
	// framework retries to call backend same command specified for parent menu
	Retry func(message *model.MessageIn, c chan<- interfaces.MessageOut)
}

func (m *MenuHandler) GetDescription() string {
	return "Handler for menu"
}

func NewMenuHandler() *MenuHandler {
	mh := MenuHandler{
		chatToMenu: map[int64]interfaces.Menu{},
		localizedPatterns: map[string]map[string]interfaces.Menu{
			defaultLang: {},
		},
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

func (m *MenuHandler) AddLocalizedMenu(item menu2.LocalizedMenu) {
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

func (m *MenuHandler) Handle(in *model.MessageIn, out interfaces.Chat) {
	menuStorage := m.getMenuStorage(in.Ctx)

	switch in.Command {
	//case menu.MenuCall:
	//	out.SendMessage(m.handleMenuCall(in, menuStorage))
	case menu2.Back:
		if currentMenu, ok := m.chatToMenu[in.Chat.ID]; ok {
			prev := currentMenu.GetPreviousMenu()
			if prev == nil {
				out.SendMessage(response.NewMessage("Nowhere to return!"))
				return
			}
			ms := m.getMenuStorage(in.Ctx)

			newPrev := ms[prev.GetCallCommand()]
			if newPrev == nil {
				in.Text = prev.GetCallCommand()
				c := make(chan interfaces.MessageOut)
				go m.Retry(in, c)
				inst := <-c

				var ok bool
				if newPrev, ok = inst.(interfaces.Menu); !ok {
					return
				}
			}
			newPrev.ForceSetMessageId(prev.GetMessageId())
			newPrev.SetChatIdIfZero(prev.GetChatId())
			newPrev.SetPreviousMenu(prev.GetPreviousMenu())
			out.SendMessage(newPrev)
		}
	default:
		out.SendMessage(m.startMenu(in, menuStorage))
	}
}

func (m *MenuHandler) startMenu(in *model.MessageIn, mStorage map[string]interfaces.Menu) interfaces.MessageOut {
	if pattern, ok := mStorage[in.Text]; ok {

		if msg, ok := m.chatToMenu[in.Chat.ID]; ok {
			pattern = pattern.GetCopy()
			pattern.ForceSetMessageId(msg.GetMessageId())
			if msg.GetCallCommand() != pattern.GetCallCommand() {
				pattern.SetPreviousMenu(msg)
			}
		} else {
			pattern.ForceSetMessageId(int64(in.MessageID))
		}
		pattern.SetChatIdIfZero(in.Chat.ID)

		return pattern
	}
	return response.NewMessage(NoMenuError + in.Text)
}

// TODO
func (m *MenuHandler) handleMenuCall(in *model.MessageIn, mStorage map[string]interfaces.Menu) interfaces.MessageOut {
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
		pattern.ForceSetMessageId(int64(in.MessageID))

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
	mn, ok := m.chatToMenu[chatId]
	if ok && mn.GetCallCommand() != menu.GetCallCommand() {
		menu.ForceSetMessageId(mn.GetMessageId())
	} else {
		menu.ForceSetMessageId(0)
	}
	m.chatToMenu[chatId] = menu
}

// TODO
func (m *MenuHandler) CanHandle(in *model.MessageIn) bool {
	// Can handler if there is a chat with menu
	if _, ok := m.chatToMenu[in.Chat.ID]; ok {
		return true
	}

	_, ok := m.localizedPatterns[defaultLang][in.Text]
	if ok {
		return true
	}

	_, ok = m.getMenuStorage(in.Ctx)[in.Text]
	return ok
}

func (m *MenuHandler) StartMenu(message *model.MessageIn, c chan<- interfaces.MessageOut) interfaces.Menu {
	ms := m.getMenuStorage(message.Ctx)
	mn, ok := ms[message.Command]
	if !ok {
		m.Retry(message, c)
		return nil
	}
	m.chatToMenu[message.Chat.ID] = mn
	return mn
}

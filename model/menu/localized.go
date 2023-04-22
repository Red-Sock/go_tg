package menu

import "github.com/RedSock/go_tg/interfaces"

type LocalizedMenu struct {
	Command    string // Printing this command calls this menu
	LangToMenu map[string]interfaces.Menu
}

func NewLocalizedMenu(command string) *LocalizedMenu {
	return &LocalizedMenu{
		Command:    command,
		LangToMenu: map[string]interfaces.Menu{},
	}
}

func (l *LocalizedMenu) AddMenu(locale string, menu interfaces.Menu) {
	l.LangToMenu[locale] = menu
}

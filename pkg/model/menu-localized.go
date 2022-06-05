package model

type LocalizedMenu struct {
	Command    string // Printing this command calls this menu
	LangToMenu map[string]Menu
}

func NewLocalizedMenu(command string) *LocalizedMenu {
	return &LocalizedMenu{
		Command:    command,
		LangToMenu: map[string]Menu{},
	}
}

func (l *LocalizedMenu) AddMenu(locale string, menu Menu) {
	l.LangToMenu[locale] = menu
}

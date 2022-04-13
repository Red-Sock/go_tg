package model

import (
	"context"
)

// Menu - menu interface
// GetName - returns name of menu (displayed at the top of the message)
// GetCallCommand - returns command (in `/model.MenuCall command` format) that instantly opens this menu
// GetTransitionCommand - returns command that (`/command format`) that replaces current open menu with this one
// GetPreviousMenu - returns previously oppened menu
// in case of LocalizedMenu realisation: methods GetName and GetPage extracts locale from context.Context
type Menu interface {
	GetName(ctx context.Context) (name string)
	GetPage(ctx context.Context, i int) (menu InlineKeyboard)
	GetCallCommand() (command string)
	GetTransitionCommand() (command string)

	GetPreviousMenu() (menu Menu)
}

// MenuPattern
// GetCallCommand - create new menu
// GetTransitionCommand - replaces current menu onto menu from command
type MenuPattern interface {
	GetCallCommand() string
	GetTransitionCommand() string
	AddPrevMenu(MenuPattern)

	Compile() Menu
}

func NewMenu(pattern MenuPattern) Menu {
	return pattern.Compile()
}

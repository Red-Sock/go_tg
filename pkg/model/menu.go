package model

import (
	"context"
)

type Menu interface {
	GetName(ctx context.Context) string
	GetPage(ctx context.Context, i int) InlineKeyboard
	GetCallCommand() string
	GetTransitionCommand() string
}

// MenuPattern
// GetCallCommand - create new menu
// GetTransitionCommand - replaces current menu onto menu from command
type MenuPattern interface {
	GetCallCommand() string
	GetTransitionCommand() string

	compile() Menu
}

func NewMenu(pattern MenuPattern) Menu {
	return pattern.compile()
}

// MenuButton - basic inline Keyboard MenuPattern item
// Name is displayed name of button
// Command - directly sends "reply message" to backend service. Message doesn't get shown in chat
type MenuButton struct {
	Name    string
	Command string
}

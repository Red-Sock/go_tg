package model

import (
	"context"
	"strconv"
)

type SimpleMenu struct {
	ChatId    int64
	MessageId int64
	Name      string
	Command   string
	Keyboard  []InlineKeyboard
}

func (m *SimpleMenu) GetPage(_ context.Context, i int) InlineKeyboard {
	if i >= len(m.Keyboard) {
		return m.Keyboard[len(m.Keyboard)-1]
	}
	if i <= 0 {
		return m.Keyboard[0]
	}
	return m.Keyboard[i]
}

func (m *SimpleMenu) GetName(_ context.Context) string {
	return m.Name
}

func (m *SimpleMenu) GetCallCommand() string {
	return MenuCall + " " + m.Name
}

func (m *SimpleMenu) GetTransitionCommand() string {
	return m.Command
}

// MenuPattern - basic inline Keyboard menu
// consists of buttons - Items (commands)
type SimpleMenuPattern struct {
	Name       string
	RawsOnPage int
	Columns    int
	EntryPoint string
	Items      []MenuButton
	keyboard   []InlineKeyboard
}

func NewMenuPattern(name string) SimpleMenuPattern {
	return SimpleMenuPattern{
		Name: name,
	}
}

func (m *SimpleMenuPattern) AddMenuButton(name, command string) {
	m.Items = append(m.Items, MenuButton{
		Name:    name,
		Command: command,
	})
}

// Adds additional command to start menu
// Usually menu gets called with MenuCall + name_of_menu
// to automate transitions between menus
// Use this to add another entry point to menu
func (m *SimpleMenuPattern) AddEntryPoint(command string) {
	m.EntryPoint = command
}

func (m *SimpleMenuPattern) GetCallCommand() string {
	return MenuCall + " " + m.Name
}

func (m *SimpleMenuPattern) GetTransitionCommand() string {
	return m.EntryPoint
}

func (m *SimpleMenuPattern) compile() Menu {
	if m.RawsOnPage == 0 {
		m.RawsOnPage = 5
	}

	pageNumber := -1

	for i, item := range m.Items {
		if i%m.RawsOnPage == 0 {
			m.keyboard = append(m.keyboard, InlineKeyboard{})
			pageNumber++
		}
		m.keyboard[pageNumber].AddButton(item.Name, item.Command)
	}

	if len(m.keyboard) > 0 {
		for i, item := range m.keyboard {
			if i > 0 {
				item.AddButton("⬅", ChangePage+strconv.Itoa(i-1))
			}

			if i < len(m.keyboard) {
				item.AddButton("➡", ChangePage+strconv.Itoa(i+1))
			}
		}
	}

	return &SimpleMenu{
		Name:     m.Name,
		Keyboard: m.keyboard,
	}
}

const (
	OpenMenu   = "/open_menu"
	MenuCall   = "/menu_call"
	ChangePage = MenuCall + " change_page"
)

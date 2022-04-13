package model

import (
	"context"
	"strconv"
)

type SimpleMenu struct {
	Name         string
	Command      string
	Keyboard     []InlineKeyboard
	previousMenu Menu
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

func (m *SimpleMenu) GetPreviousMenu() Menu {
	return m.previousMenu
}

// SimpleMenuPattern - basic inline Keyboard menu
// consists of buttons - Items (commands)
type SimpleMenuPattern struct {
	Name       string
	RawsOnPage int
	Columns    int
	EntryPoint string
	keyboard   InlineKeyboard

	previousMenu MenuPattern
}

func NewMenuPattern(name string) SimpleMenuPattern {
	return SimpleMenuPattern{
		Name: name,
	}
}

// AddButton - adds button with displayed name, that executes command when pressed
func (m *SimpleMenuPattern) AddButton(name, command string) {
	m.keyboard.AddButton(name, command)
}

// AddStandAloneButton - adds button (like AddButton) but puts it between rows
// like this
// [button1][button2]
// [standalone button] <-- added with AddStandAloneButton method
// [button3][button4]
func (m *SimpleMenuPattern) AddStandAloneButton(text, value string) {
	m.keyboard.AddStandAloneButton(text, value)
}

// AddButtonWithMenu - adds button with displayed name, that opens given menu
func (m *SimpleMenuPattern) AddButtonWithMenu(name string, menu MenuPattern) {
	menu.AddPrevMenu(m)
	m.keyboard.AddButton(name, menu.GetCallCommand())
}

// AddPrevMenu - adds previous menu
func (m *SimpleMenuPattern) AddPrevMenu(pattern MenuPattern) {
	m.previousMenu = pattern
}

// AddEntryPoint - adds additional command to start menu
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

func (m *SimpleMenuPattern) Compile() Menu {
	if m.RawsOnPage == 0 {
		m.RawsOnPage = 5
	}
	var keyboard []InlineKeyboard

	if len(m.keyboard.btns) > m.RawsOnPage {
		for i := 0; i < len(m.keyboard.btns); i += m.RawsOnPage {
			keyboard = append(keyboard, InlineKeyboard{btns: m.keyboard.btns[i:min(len(m.keyboard.btns), i+m.RawsOnPage)]})
			//if item.StandAlone {
			//	keyboard[pageNumber].AddStandAloneButton(item.Text, item.Value)
			//} else {
			//	keyboard[pageNumber].AddButton(item.Text, item.Value)
			//}
		}

		if len(keyboard) > 0 {
			for i, item := range keyboard {
				if i > 0 {
					item.AddButton("⬅", ChangePage+strconv.Itoa(i-1))
				}

				if i < len(keyboard) {
					item.AddButton("➡", ChangePage+strconv.Itoa(i+1))
				}
			}
		}
	} else {
		keyboard = append(keyboard, m.keyboard)
	}

	return &SimpleMenu{
		Name:     m.Name,
		Keyboard: keyboard,
	}
}

const (
	OpenMenu   = "/open_menu"
	MenuCall   = "/menu_call"
	ChangePage = MenuCall + " change_page"
	GoBack     = MenuCall + " back"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

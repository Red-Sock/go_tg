package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type SimpleMenu struct {
	command       string
	displayedName string

	columnsPerPage uint8
	rowsPerPage    uint8

	buttons       InlineKeyboard
	pagesPatterns []InlineKeyboard
	pages         []*tgbotapi.InlineKeyboardMarkup

	previousMenu Menu
}

func NewSimpleMenu(command, name string) *SimpleMenu {
	return &SimpleMenu{
		command:        command,
		displayedName:  name,
		columnsPerPage: ColumnsDefaultAmount,
		rowsPerPage:    RowsDefaultAmount,
	}
}

func (m *SimpleMenu) GetName() (name string) {
	return m.displayedName
}

func (m *SimpleMenu) GetCallCommand() (command string) {
	return MenuCall + " " + m.command
}

func (m *SimpleMenu) GetPage(numbers ...int) (menu *tgbotapi.InlineKeyboardMarkup) {
	page := 0
	if len(numbers) != 0 {
		page = numbers[0]
	}
	return m.pages[page]
}

func (m *SimpleMenu) SetPreviousMenu(menu Menu) {
	m.previousMenu = menu
}

func (m *SimpleMenu) GetPreviousMenu() (menu Menu) {
	return m.previousMenu
}

// SetColumnsAmount - sets amount of columns per page
func (m *SimpleMenu) SetColumnsAmount(col uint8) {
	m.columnsPerPage = col
	m.rebuild()
}

// SetRowsAmount - sets amount of rows displayed per page
func (m *SimpleMenu) SetRowsAmount(rows uint8) {
	m.rowsPerPage = rows
	m.rebuild()
}

// AddButton - adds button to the end of menu
// Careful!! On every call rebuilds whole menu keyboard
// More preferred way to create menu is to use AddPage method
func (m *SimpleMenu) AddButton(name, command string) {

	m.buttons.AddButton(name, command)
	m.rebuild()
}

// AddStandAloneButton - on every call rebuilds whole menu keyboard
// Careful!! On every call rebuilds whole menu keyboard
// More preferred way to create menu is to use AddPage method
func (m *SimpleMenu) AddStandAloneButton(name, command string) {
	m.buttons.AddStandAloneButton(name, command)
	m.rebuild()
}

func (m *SimpleMenu) AddPage(menu InlineKeyboard) {
	m.pagesPatterns = append(m.pagesPatterns, menu)
}

func (m *SimpleMenu) rebuild() {
	var curPage, curRows, curCol uint8

	menu := make([]InlineKeyboard, 1)

	for _, item := range m.buttons.btns {
		menu[curPage].btns = append(menu[curPage].btns, item)

		if item.IsStandAlone {
			curRows++
			curCol = 0
		} else {
			curCol++
			if curCol >= m.columnsPerPage {
				curRows++
			}
		}

		if curRows >= m.rowsPerPage {
			curPage++
			curRows = 0
			menu = append(menu, InlineKeyboard{})
		}
	}
	m.pagesPatterns = menu
}

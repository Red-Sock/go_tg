package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/AlexSkilled/go_tg/interfaces"
)

// SimpleMenu is a basic menu with chat buttons
// If both ChatId And MessageId specified - replaces message with this one
type SimpleMenu struct {
	ChatId        int64
	MessageId     int64
	command       string
	displayedName string

	columnsPerPage uint8
	rowsPerPage    uint8

	buttons       InlineKeyboard
	pagesPatterns []InlineKeyboard
	pages         []*tgbotapi.InlineKeyboardMarkup

	previousMenu interfaces.Menu
}

func (m *SimpleMenu) GetCopy() interfaces.Menu {
	mCopy := *m
	return &mCopy
}

func NewSimple(name, command string) *SimpleMenu {
	return &SimpleMenu{
		command:        command,
		displayedName:  name,
		columnsPerPage: ColumnsDefaultAmount,
		rowsPerPage:    RowsDefaultAmount,
		pages:          make([]*tgbotapi.InlineKeyboardMarkup, 0, 1),
	}
}

func (m *SimpleMenu) GetName() (name string) {
	return m.displayedName
}

// GetPage returns compiled page if nessesary
func (m *SimpleMenu) GetPage(numbers ...int) (menu *tgbotapi.InlineKeyboardMarkup) {
	page := 0
	if len(numbers) != 0 {
		page = numbers[0]
	}
	return m.pages[page]
}

func (m *SimpleMenu) GetCallCommand() string {
	return m.command
}

func (m *SimpleMenu) SetPreviousMenu(menu interfaces.Menu) {
	m.previousMenu = menu
	m.rebuild()
}

func (m *SimpleMenu) GetPreviousMenu() (menu interfaces.Menu) {
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
	m.pages = append(m.pages, menu.ToMarkup())
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
	if m.previousMenu != nil {
		menu[0].btns = append(menu[0].btns, button{
			Text:         "<<",
			Value:        "/back",
			IsStandAlone: true,
		})
	}

	m.pagesPatterns = menu
	m.pages = make([]*tgbotapi.InlineKeyboardMarkup, 0, len(menu))

	for _, item := range m.pagesPatterns {
		item.Columns = m.columnsPerPage
		item.Rows = m.rowsPerPage
		m.pages = append(m.pages, item.ToMarkup())
	}
}

func (m *SimpleMenu) SetChatIdIfZero(c int64) {
	if m.ChatId == 0 {
		m.ChatId = c
	}
}

func (m *SimpleMenu) GetMessage() tgbotapi.Chattable {
	if m.MessageId != 0 {
		return tgbotapi.NewEditMessageTextAndMarkup(m.ChatId, int(m.MessageId), m.displayedName, *m.pages[0])
	}
	message := tgbotapi.NewMessage(m.ChatId, m.displayedName)

	message.DisableWebPagePreview = true
	message.ReplyMarkup = m.pages[0]

	return message
}

func (m *SimpleMenu) GetChatId() int64 {
	return m.ChatId
}

func (m *SimpleMenu) ForceSetMessageId(id int64) {
	m.MessageId = id
}

func (m *SimpleMenu) GetMessageId() int64 {
	return m.MessageId
}

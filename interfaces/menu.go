package interfaces

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Menu - menu interface
// GetName - returns name of menu (displayed at the top of the message)
// GetPage - returns first page if number is not specified. If specified - return buttons on page number. For more details
// look at the implementation e.g. SimpleMenu
// GetCallCommand - returns command (in `/model.MenuCall command` format) that instantly opens this menu
// GetTransitionCommand - returns command that (`/command format`) that replaces currently opened menu with this one
// GetPreviousMenu - returns previously opened menu
type Menu interface {
	MessageOut
	GetName() (name string)
	GetPage(number ...int) (menu *tgbotapi.InlineKeyboardMarkup)
	GetCallCommand() (command string)

	SetPreviousMenu(menu Menu)
	GetPreviousMenu() (menu Menu)

	GetCopy() Menu
}

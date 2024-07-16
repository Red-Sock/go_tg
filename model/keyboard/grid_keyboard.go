package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GridKeyboard struct {
	buttons []Button

	Columns uint8
	Rows    uint8

	isReplyKeyboard
}

func (b *GridKeyboard) AddButton(button Button) {
	b.buttons = append(b.buttons, button)
}

func (b *GridKeyboard) AddStandAloneButton(text, value string) {
	b.buttons = append(b.buttons, Button{
		Text:         text,
		Value:        value,
		IsStandAlone: true,
	})
}

func (b *GridKeyboard) ToInlineMarkup() *tgbotapi.InlineKeyboardMarkup {
	if b.Columns == 0 {
		b.Columns = ColumnsDefaultAmount
	}

	if b.Rows == 0 {
		b.Rows = RowsDefaultAmount
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 1, b.Rows)

	var cRaw, cCol uint8

	processedButtons := 0
	for _, btn := range b.buttons {
		btnMark := tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Value)

		if btn.IsStandAlone {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{btnMark})
			cRaw++
			continue
		}

		if cCol >= b.Columns {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{btnMark})
			cRaw++
			cCol = 1
		} else {
			rows[cRaw] = append(rows[cRaw], btnMark)
			cCol++
		}

		processedButtons++

		if cRaw >= b.Rows {
			return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
		}
	}
	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func (b *GridKeyboard) ToKeyboardMarkup() *tgbotapi.ReplyKeyboardMarkup {
	if b.Columns == 0 {
		b.Columns = ColumnsDefaultAmount
	}

	if b.Rows == 0 {
		b.Rows = RowsDefaultAmount
	}

	rows := make([][]tgbotapi.KeyboardButton, 1, b.Rows)

	var cRaw, cCol uint8

	processedButtons := 0
	for _, btn := range b.buttons {
		btnMark := tgbotapi.NewKeyboardButton(btn.Text)

		if btn.IsStandAlone {
			rows = append(rows, []tgbotapi.KeyboardButton{btnMark})
			cRaw++
			continue
		}

		if cCol >= b.Columns {
			rows = append(rows, []tgbotapi.KeyboardButton{btnMark})
			cRaw++
			cCol = 1
		} else {
			rows[cRaw] = append(rows[cRaw], btnMark)
			cCol++
		}

		processedButtons++

		if cRaw >= b.Rows {
			return &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows}
		}
	}
	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows}
}

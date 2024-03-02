package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboard struct {
	btns []button

	Columns uint8
	Rows    uint8

	IsReplyKeyboard bool
}

type button struct {
	Text         string
	Value        string
	IsStandAlone bool
}

func (b *Keyboard) AddButton(text, value string) {
	b.btns = append(b.btns, button{
		Text:  text,
		Value: value,
	})
}

func (b *Keyboard) AddStandAloneButton(text, value string) {
	b.btns = append(b.btns, button{
		Text:         text,
		Value:        value,
		IsStandAlone: true,
	})
}

func (b *Keyboard) ToInlineMarkup() *tgbotapi.InlineKeyboardMarkup {
	if b.Columns == 0 {
		b.Columns = ColumnsDefaultAmount
	}

	if b.Rows == 0 {
		b.Rows = RowsDefaultAmount
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 1, b.Rows)

	var cRaw, cCol uint8

	processedButtons := 0
	for _, btn := range b.btns {
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

func (b *Keyboard) ToKeyboardMarkup() *tgbotapi.ReplyKeyboardMarkup {
	if b.Columns == 0 {
		b.Columns = ColumnsDefaultAmount
	}

	if b.Rows == 0 {
		b.Rows = RowsDefaultAmount
	}

	rows := make([][]tgbotapi.KeyboardButton, 1, b.Rows)

	var cRaw, cCol uint8

	processedButtons := 0
	for _, btn := range b.btns {
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

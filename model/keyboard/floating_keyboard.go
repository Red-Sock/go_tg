package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FloatingKeyboard struct {
	rows [][]Button

	Columns uint8
	Rows    uint8

	isReplyKeyboard
}

func (b *FloatingKeyboard) AddRow(row []Button) {
	b.rows = append(b.rows, row)
}

func (b *FloatingKeyboard) ToInlineMarkup() *tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, b.Rows)

	for i, row := range b.rows {
		rows = append(rows, make([]tgbotapi.InlineKeyboardButton, 0, len(row)))
		for _, btn := range row {
			rows[i] = append(rows[i], tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Value))
		}
	}
	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func (b *FloatingKeyboard) ToKeyboardMarkup() *tgbotapi.ReplyKeyboardMarkup {

	rows := make([][]tgbotapi.KeyboardButton, 0, b.Rows)

	for i, row := range b.rows {
		rows = append(rows, make([]tgbotapi.KeyboardButton, 0, len(row)))
		for _, btn := range row {
			rows[i] = append(rows[i], tgbotapi.NewKeyboardButton(btn.Text))
		}
	}
	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows}
}

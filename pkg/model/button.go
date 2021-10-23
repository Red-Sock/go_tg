package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type button struct {
	Text  string
	Value string
}

type InlineKeyboard struct {
	btns []button

	Columns int
}

func (b *InlineKeyboard) AddButton(text, value string) {
	b.btns = append(b.btns, button{
		Text:  text,
		Value: value,
	})
}
func (b *InlineKeyboard) ToMarkup() *tgbotapi.InlineKeyboardMarkup {
	if b.Columns == 0 {
		b.Columns = 1
	}

	finalButtonsSet := make([][]tgbotapi.InlineKeyboardButton, 0, 1)
	i := 0
	raw := -1
	for i < len(b.btns) {
		if i%b.Columns == 0 {
			finalButtonsSet = append(finalButtonsSet, make([]tgbotapi.InlineKeyboardButton, 0, b.Columns))
			raw++
		}
		finalButtonsSet[raw] = append(finalButtonsSet[raw],
			tgbotapi.NewInlineKeyboardButtonData(b.btns[i].Text, b.btns[i].Value))
		i++
	}

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: finalButtonsSet}
}

type Keyboard struct {
	btns []string

	Columns int

	ResizeKeyboard bool `json:"resize_keyboard"`
}

func (b *Keyboard) AddKey(value string) {
	b.btns = append(b.btns, value)
}
func (b *Keyboard) toMarkup() *tgbotapi.ReplyKeyboardMarkup {
	if b.Columns == 0 {
		b.Columns = 1
	}

	finalButtonsSet := make([][]tgbotapi.KeyboardButton, 0, 1)
	i := 0
	raw := -1
	for i < len(b.btns) {
		if i%b.Columns == 0 {
			finalButtonsSet = append(finalButtonsSet, make([]tgbotapi.KeyboardButton, 0, b.Columns))
			raw++
		}
		finalButtonsSet[raw] = append(finalButtonsSet[raw],
			tgbotapi.NewKeyboardButton(b.btns[i]))
		i++
	}

	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: finalButtonsSet, ResizeKeyboard: true}
}

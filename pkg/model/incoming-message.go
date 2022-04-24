package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageIn struct {
	Command string
	Args    []string

	*tgbotapi.Message
	m Message
}

type Message struct {
}

func (m *Message) Delete() {

}

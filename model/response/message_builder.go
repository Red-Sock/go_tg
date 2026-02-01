package response

import (
	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model/keyboard"
)

type Builder interface {
	SetText(text string) Builder
	SetKeyboard(keyboard keyboard.Keyboard) Builder

	EditMessage(messageId int) Builder
	Delete(messageId int) Builder

	Build() interfaces.MessageOut
}
type msgBuilderFunc func(builder *responseBuilder) interfaces.MessageOut

type responseBuilder struct {
	text     string
	keyboard keyboard.Keyboard

	messageId *int

	buildFunc msgBuilderFunc
}

func New() Builder {
	return &responseBuilder{
		buildFunc: buildSimpleMessage,
	}
}

func (r *responseBuilder) SetText(text string) Builder {
	r.text = text
	return r
}

func (r *responseBuilder) SetKeyboard(keyboard keyboard.Keyboard) Builder {
	r.keyboard = keyboard
	return r
}

func (r *responseBuilder) EditMessage(messageId int) Builder {
	r.messageId = &messageId
	r.buildFunc = buildEditMessage
	return r
}

func (r *responseBuilder) Delete(messageId int) Builder {
	r.messageId = &messageId
	r.buildFunc = buildDeleteMessage
	return r
}

func (r *responseBuilder) Build() interfaces.MessageOut {
	return r.buildFunc(r)
}

func buildSimpleMessage(builder *responseBuilder) interfaces.MessageOut {
	m := &MessageOut{
		Text: builder.text,
	}

	m.AddKeyboard(builder.keyboard)

	return m
}

func buildEditMessage(builder *responseBuilder) interfaces.MessageOut {
	if builder.messageId == nil {
		return buildSimpleMessage(builder)
	}

	m := &EditMessage{
		Text:      builder.text,
		Keys:      builder.keyboard,
		MessageId: int64(*builder.messageId),
	}

	return m
}

func buildDeleteMessage(builder *responseBuilder) interfaces.MessageOut {
	if builder.messageId == nil {
		return buildSimpleMessage(builder)
	}

	m := &DeleteMessage{
		MessageId: int64(*builder.messageId),
	}

	return m
}

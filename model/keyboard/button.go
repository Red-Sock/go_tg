package keyboard

type Button struct {
	Text         string
	Value        string
	IsStandAlone bool
}

func NewButton(text, value string) Button {
	return Button{
		Text:         text,
		Value:        value,
		IsStandAlone: false,
	}
}

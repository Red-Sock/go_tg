package keyboard

type Builder interface {
	AddButton(b Button) Builder
}

type keyboardBuilder struct {
	buttons []Button
}

func New() Builder {
	return &keyboardBuilder{}
}

func (k *keyboardBuilder) AddButton(b Button) Builder {
	k.buttons = append(k.buttons, b)

	return k
}

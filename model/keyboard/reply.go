package keyboard

type isReplyKeyboard bool

func (i *isReplyKeyboard) IsReplyKeyboard() bool {
	return bool(*i)
}

func (i *isReplyKeyboard) SetIsReplyKeyboard(v bool) {
	*i = isReplyKeyboard(v)
}

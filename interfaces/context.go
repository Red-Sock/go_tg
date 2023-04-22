package interfaces

import (
	"context"

	"github.com/RedSock/go_tg/model"
)

type ExternalContext interface {
	GetContext(in *model.MessageIn) (context.Context, error)
}

type GetContextFunc func(in *model.MessageIn) (context.Context, error)

func (g GetContextFunc) GetContext(in *model.MessageIn) (context.Context, error) {
	return g(in)
}

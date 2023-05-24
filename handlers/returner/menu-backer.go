package returner

import (
	"context"
	"strings"

	"github.com/Red-Sock/go_tg/handlers/commands"
	"github.com/Red-Sock/go_tg/interfaces"
	"github.com/Red-Sock/go_tg/model"
	"github.com/Red-Sock/go_tg/model/response"
)

type Handler struct {
	Handlers map[string]interfaces.CommandHandler
}

func (h *Handler) Handle(in *model.MessageIn, out interfaces.Chat) {
	if len(in.Args) < 1 {
		out.SendMessage(&response.MessageOut{Text: "Cannot return. Need 1 argument {{command to return to}}. E.g. /return /main_menu"})
		return
	}
	prox := &chatProxy{
		realChat: out,
	}

	handler, ok := h.Handlers[in.Args[0]]
	if !ok {
		out.SendMessage(&response.MessageOut{Text: "Cannot return to " + in.Args[0] + ". Unknown handler"})
		return
	}
	mesIn := &model.MessageIn{
		Args:    in.Args[1:],
		Ctx:     in.Ctx,
		Message: in.Message,
	}
	mesIn.Text = strings.Join(in.Args, " ")

	handler.Handle(mesIn, prox)

	prox.savedMessage.ForceSetMessageId(int64(in.MessageID))

	out.SendMessage(prox.savedMessage)
}

func (h *Handler) GetDescription() string {
	return "returns to menu by given command "
}

func (h *Handler) GetCommand() string {
	return commands.Return
}

type chatProxy struct {
	savedMessage interfaces.MessageOut
	realChat     interfaces.Chat
}

func (c *chatProxy) GetInput(ctx context.Context) (*model.MessageIn, error) {
	return c.realChat.GetInput(ctx)
}
func (c *chatProxy) SendMessage(out interfaces.MessageOut) {
	c.savedMessage = out
}

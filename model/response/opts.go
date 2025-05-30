package response

import (
	"github.com/Red-Sock/go_tg/model/media"
)

type opt func(m *MessageOut)

func WithMedia(mediaFiles ...media.Media) opt {
	return func(m *MessageOut) {
		m.Media = append(m.Media, mediaFiles...)
	}
}

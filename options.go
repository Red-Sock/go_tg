package go_tg

import (
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type opt func(bot *Bot)

func WithLogger(logger logrus.FieldLogger) opt {
	return func(bot *Bot) {
		bot.logger = logger
	}
}

// WithOnlyDirectCalls - set value for "Direct calls" flag
// When "Direct calls" flag is set in group chats bot will
// be trigger only when command is called with bot tag
//
//	e.g.
//		with onlyDirectCalls = true
//		command /version - won't be handled by bot
//		but command /verstion@bot_name_here_bot - will trigger handler
func WithOnlyDirectCalls(v bool) opt {
	return func(bot *Bot) {
		bot.onlyDirectCalls = v
	}
}

// WithHTTPClient replaces the HTTP client used for all Telegram API calls.
// Useful for custom TLS config, timeouts, or bringing your own proxy transport.
func WithHTTPClient(client tgbotapi.HTTPClient) opt {
	return func(bot *Bot) {
		bot.httpClient = client
	}
}

// WithProxy sets an HTTP/HTTPS proxy for all Telegram API calls.
// rawURL must be a valid URL (e.g. "http://proxy.example.com:8080").
// Panics on a malformed URL — treat this as a configuration error.
func WithProxy(rawURL string) opt {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic("go_tg: invalid proxy URL: " + err.Error())
	}
	return func(bot *Bot) {
		bot.httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}
	}
}

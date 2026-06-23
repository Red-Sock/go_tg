package go_tg

import (
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type loggingHTTPClient struct {
	inner  tgbotapi.HTTPClient
	logger logrus.FieldLogger
}

func (c *loggingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.logger.Debugf("tg http %s %s", req.Method, req.URL)
	return c.inner.Do(req)
}

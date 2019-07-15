package discord

import (
	"net/http"
)

type (
	client struct {
		Endpoint   string
		HTTPClient *http.Client
		Token      string
	}

	BotClient interface {

	}
)

func NewClient(token string) *client {
	return &client{
		Endpoint: "https://discordapp.com/api/v6",
		HTTPClient: &http.Client{},
		Token: token,
	}
}

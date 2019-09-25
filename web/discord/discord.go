package discord

import (
	"io/ioutil"
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type (
	bot struct {
		URL string
		HTTPClient *http.Client
		Token string
	}

	BotClient interface {
		GetGuilds() (guild *[]Guild, err error)
		GetEmojis(guildid string) (*[]Emoji, error)
		GetEmoji(guildid string, emojiid string) (*Emoji, error)
	}

	user struct {
		URL string
		Config     *oauth2.Config
	}

	UserClient interface {
		GetGuilds(token *oauth2.Token) (guild *[]Guild, err error)
	}
)

func NewBotClient(token string) *bot {
	return &bot{
		URL: "https://discordapp.com/api/v6",
		HTTPClient: &http.Client{},
		Token: token,
	}
}

func NewUserClient(config *oauth2.Config) *user {
	return &user{
		URL: "https://discordapp.com/api/v6",
		Config: config,
	}
}

func (b *bot) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", b.URL + endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bot " + b.Token)

	res, err := b.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return parseResponse(res)
}

func (u *user) get(endpoint string, token *oauth2.Token) ([]byte, error) {
	client := u.Config.Client(context.Background(), token)

	res, err := client.Get(u.URL + endpoint)
	if err != nil {
		return nil, err
	}
	
	return parseResponse(res)
}

func parseResponse(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}
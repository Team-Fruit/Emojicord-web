package discord

import (
	"encoding/json"

	"golang.org/x/oauth2"
)

type Guild struct {
	Owner       bool   `json:"owner"`
	Permissions int    `json:"permissions"`
	Icon        string `json:"icon"`
	ID          string `json:"id"`
	Name        string `json:"name"`
}

func (u *user) GetGuilds(token *oauth2.Token) (*[]Guild, error) {
	body, err := u.get("/users/@me/guilds", token)
	if err != nil {
		return nil, err
	}

	var guilds []Guild
	if err = json.Unmarshal(body, &guilds); err != nil {
		return nil, err
	}

	return &guilds, err
}

func (b *bot) GetGuilds() (*[]Guild, error) {
	body, err := b.get("/users/@me/guilds")
	if err != nil {
		return nil, err
	}

	var guilds []Guild
	if err = json.Unmarshal(body, &guilds); err != nil {
		return nil, err
	}

	if len(guilds) < 100 {
		return &guilds, nil
	}

	for {
		last := guilds[len(guilds)-1].ID

		body, err := b.get("/users/@me/guilds?after=" + last)
		if err != nil {
			return nil, err
		}
	
		var page []Guild
		if err = json.Unmarshal(body, &page); err != nil {
			return nil, err
		}
	
		guilds = append(guilds, page...)

		if len(page) < 100 {
			break
		}
	}

	return &guilds, err
}
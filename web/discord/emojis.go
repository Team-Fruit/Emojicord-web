package discord

import (
	"encoding/json"
)

type Emoji struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	User     EmojiUser `json:"user"`
	Animated bool      `json:"animated"`
}

type EmojiUser struct {
	UserName      string `json:"username"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
	avatar        string `json:"avatar"`
}

func (b *bot) GetEmojis(guildid string) (*[]Emoji, error) {
	body, err := b.get("/guilds/" + guildid + "/emojis")
	if err != nil {
		return nil, err
	}

	var emojis []Emoji
	if err = json.Unmarshal(body, &emojis); err != nil {
		return nil, err
	}

	return &emojis, nil
}

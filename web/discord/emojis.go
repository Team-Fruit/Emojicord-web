package discord

import (
	"encoding/json"
)

type (
	EmojiUser struct {
		UserName      string `json:"username"`
		Discriminator string `json:"discriminator"`
		ID            string `json:"id"`
		avatar        string `json:"avatar"`
	}

	Emoji struct {
		ID       string    `json:"id"`
		Name     string    `json:"name"`
		User     EmojiUser `json:"user"`
		Animated bool      `json:"animated"`
		GuildID  string
	}
)

func (b *bot) GetEmojis(guildid string) ([]*Emoji, error) {
	body, err := b.get("/guilds/" + guildid + "/emojis")
	if err != nil {
		return nil, err
	}

	var emojis []*Emoji
	if err = json.Unmarshal(body, &emojis); err != nil {
		return nil, err
	}

	for i := range emojis {
		emojis[i].GuildID = guildid
	}

	return emojis, nil
}

func (b *bot) GetEmoji(guildid string, emojiid string) (*Emoji, error) {
	body, err := b.get("/guilds/"+guildid+"/emojis/"+emojiid)
	if err != nil {
		return nil, err
	}

	var emoji Emoji
	if err = json.Unmarshal(body, &emoji); err != nil {
		return nil, err
	}

	emoji.GuildID = guildid

	return &emoji, nil
}

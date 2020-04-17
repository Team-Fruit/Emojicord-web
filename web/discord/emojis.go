package discord

import (
	"encoding/json"
	"strconv"
)

type (
	EmojiUser struct {
		UserName      string `json:"username"`
		Discriminator string `json:"discriminator"`
		ID            uint64 `json:"id,string"`
		Avatar        string `json:"avatar"`
	}

	Emoji struct {
		ID       uint64    `json:"id,string"`
		Name     string    `json:"name"`
		User     EmojiUser `json:"user"`
		Animated bool      `json:"animated"`
		GuildID  uint64
	}
)

func (b *bot) GetEmojis(guildid uint64) ([]*Emoji, error) {
	body, err := b.get("/guilds/" + strconv.FormatUint(guildid, 10) + "/emojis")
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

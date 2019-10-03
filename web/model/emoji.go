package model

import (
	"github.com/Team-Fruit/Emojicord-web/web/discord"
)

type (
	Emoji struct {
		ID                string `json:"id" db:"id"`
		GuildID           string `json:"guildid" db:"guild_id"`
		Name              string `json:"name" db:"name"`
		Animated          bool   `json:"animated" db:"is_animated"`
		UserID            string `json:"userid" db:"user_id"`
		UserName          string `json:"username" db:"user_username"`
		UserDiscriminator string `json:"userdiscriminator" db:"user_discriminator"`
		UserAvatar        string `json:"useravatar" db:"user_avatar"`
		Enabled           bool   `json:"enabled" db:"is_enabled"`
	}
)

func (m *model) AddEmojis(emojis []*discord.Emoji) (err error) {
	tx, err := m.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for i := range emojis {
		_, err = tx.Exec(`INSERT INTO discord_emojis
						VALUES (?, ?, ?, ?, ?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						user_id = VALUES(user_id),
						name = VALUES(name)`,
			emojis[i].ID,
			emojis[i].GuildID,
			emojis[i].Name,
			emojis[i].Animated,
			emojis[i].User.ID,
			emojis[i].User.UserName,
			emojis[i].User.Discriminator,
			emojis[i].User.Avatar)
		if err != nil {
			return
		}
	}

	return
}

func (m *model) AddEmoji(emoji *discord.Emoji) (err error) {
	_, err = m.db.Exec(`INSERT INTO discord_emojis
						VALUES (?, ?, ?, ?, ?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						user_id = VALUES(user_id),
						name = VALUES(name)`,
		emoji.ID,
		emoji.GuildID,
		emoji.Name,
		emoji.Animated,
		emoji.User.ID,
		emoji.User.UserName,
		emoji.User.Discriminator,
		emoji.User.Avatar)
	return
}

func (m *model) AddUserEmojis(userid string) (err error) {
	_, err = m.db.Exec(`INSERT IGNORE INTO users__discord_emojis (user_id, emoji_id, is_enabled) 
						SELECT users__discord_guilds.user_id, discord_emojis.id, true FROM discord_emojis 
						INNER JOIN users__discord_guilds 
						ON discord_emojis.guild_id = users__discord_guilds.guild_id 
						AND users__discord_guilds.user_id = ?`,
		userid)
	return
}

func (m *model) GetUserEmojis(userid string) ([]*Emoji, error) {
	emojis := []*Emoji{}
	if err := m.db.Select(&emojis, `SELECT discord_emojis.id, 
									discord_emojis.guild_id, 
									discord_emojis.name, 
									discord_emojis.is_animated, 
									discord_emojis.user_id, 
									discord_emojis.user_username, 
									discord_emojis.user_discriminator, 
									discord_emojis.user_avatar, 
									users__discord_emojis.is_enabled 
									FROM discord_emojis 
									INNER JOIN users__discord_emojis 
									ON discord_emojis.id = users__discord_emojis.emoji_id 
									WHERE users__discord_emojis.user_id = ?;`, userid); err != nil {
		return nil, err
	}
	return emojis, nil
}

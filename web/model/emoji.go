package model

import (
	"github.com/Team-Fruit/Emojicord-web/web/discord"
)

type (
	Emoji struct {
		ID       string `db:"id"`
		GuildID  string `db:"guild_id"`
		UserID   string `db:"user_id"`
		Name     string `db:"name"`
		Animated bool   `db:"is_animated"`
	}

	UserEmoji struct {
		UserID  string `db:"user_id"`
		EmojiID string `db:"emoji_id"`
		Enabled bool   `db:"is_enabled"`
	}
)

func (m *model) AddEmojis(emojis *[]discord.Emoji) (err error) {
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

	for _, emoji := range *emojis {
		_, err = tx.Exec(`INSERT INTO discord_emojis
						VALUES (?, ?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						user_id = VALUES(user_id),
						name = VALUES(name)`,
			emoji.ID,
			emoji.GuildID,
			emoji.User.ID,
			emoji.Name,
			emoji.Animated)
		if err != nil {
			return
		}
	}

	return
}

func (m *model) AddEmoji(emoji *discord.Emoji) (err error) {
	_, err = m.db.Exec(`INSERT INTO discord_emojis
						VALUES (?, ?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						user_id = VALUES(user_id),
						name = VALUES(name)`,
		emoji.ID,
		emoji.GuildID,
		emoji.User.ID,
		emoji.Name,
		emoji.Animated)
	return
}

func (m *model) AddUserEmojis(userEmojis *[]UserEmoji) (err error) {
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

	for _, ue := range *userEmojis {
		_, err = tx.Exec(`INSERT INTO users__discord_emojis
						VALUES (?, ?, ?)
						ON DUPLICATE KEY UPDATE
						is_enabled = VALUES(is_enabled)`,
			ue.UserID,
			ue.EmojiID,
			ue.Enabled)
		if err != nil {
			return
		}
	}

	return
}

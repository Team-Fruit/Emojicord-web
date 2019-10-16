package model

import (
	"github.com/jmoiron/sqlx"

	"github.com/Team-Fruit/Emojicord-web/web/discord"
)

type (
	Emoji struct {
		ID       string `json:"id" db:"id"`
		GuildID  string `json:"guildid" db:"guild_id"`
		Name     string `json:"name" db:"name"`
		Animated bool   `json:"animated" db:"is_animated"`
		UserID   string `json:"userid" db:"user_id"`
		Enabled  bool   `json:"enabled" db:"is_enabled"`
	}

	EmojiUser struct {
		UserID            string `json:"id" db:"id"`
		UserName          string `json:"name" db:"username"`
		UserDiscriminator string `json:"discriminator" db:"discriminator"`
		UserAvatar        string `json:"avatar" db:"avatar"`
	}

	UpdateEmoji struct {
		UserID  string
		EmojiID string
		Enabled  bool
	}

	UpdateEmojis struct {
		UserID  string
		EmojiID []string
		Enabled  bool
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
		_, err = tx.Exec(`INSERT INTO discord_emojis_users
						VALUES (?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						username = VALUES(username),
						discriminator = VALUES(discriminator),
						avatar = VALUES(avatar)`,
			emojis[i].User.ID,
			emojis[i].User.UserName,
			emojis[i].User.Discriminator,
			emojis[i].User.Avatar)
		if err != nil {
			return
		}

		_, err = tx.Exec(`INSERT INTO discord_emojis
						VALUES (?, ?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						user_id = VALUES(user_id),
						name = VALUES(name)`,
			emojis[i].ID,
			emojis[i].GuildID,
			emojis[i].Name,
			emojis[i].Animated,
			emojis[i].User.ID)
		if err != nil {
			return
		}
	}

	return
}

func (m *model) AddEmoji(emoji *discord.Emoji) (err error) {
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

	_, err = tx.Exec(`INSERT INTO discord_emojis_users
						VALUES (?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						username = VALUES(username),
						discriminator = VALUES(discriminator),
						avatar = VALUES(avatar)`,
		emoji.User.ID,
		emoji.User.UserName,
		emoji.User.Discriminator,
		emoji.User.Avatar)
	if err != nil {
		return
	}

	_, err = tx.Exec(`INSERT INTO discord_emojis
						VALUES (?, ?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						user_id = VALUES(user_id),
						name = VALUES(name)`,
		emoji.ID,
		emoji.GuildID,
		emoji.Name,
		emoji.Animated,
		emoji.User.ID)
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
									users__discord_emojis.is_enabled 
									FROM discord_emojis 
									INNER JOIN users__discord_emojis 
									ON discord_emojis.id = users__discord_emojis.emoji_id 
									WHERE users__discord_emojis.user_id = ?`, userid); err != nil {
		return nil, err
	}
	return emojis, nil
}

func (m *model) GetEmojiUsers(userid string) ([]*EmojiUser, error) {
	users := []*EmojiUser{}
	if err := m.db.Select(&users, `SELECT * FROM discord_emojis_users
									WHERE id IN 
									(SELECT discord_emojis.user_id FROM discord_emojis
									INNER JOIN users__discord_emojis
									ON discord_emojis.id = users__discord_emojis.emoji_id
									WHERE users__discord_emojis.user_id = ?)`, userid); err != nil {
		return nil, err
	}
	return users, nil
}

func (m *model) UpdateUserEmoji(obj UpdateEmoji) (err error) {
	_, err = m.db.NamedExec(`UPDATE users__discord_emojis SET is_enabled = :enabled
							WHERE user_id = :userid AND emoji_id = :emojiid`, obj)
	return
}

func (m *model) UpdateUserEmojis(obj UpdateEmojis) (err error) {
	query, args, err := sqlx.Named(`UPDATE users__discord_emojis SET is_enabled = :enabled 
									WHERE user_id = :userid AND emoji_id IN (:emojiid)`, obj)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = m.db.Rebind(query)
	_, err = m.db.Exec(query, args...)
	return
}

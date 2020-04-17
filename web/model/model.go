package model

import (
	"github.com/jmoiron/sqlx"

	"github.com/Team-Fruit/Emojicord-web/web/discord"
)


type (
	model struct {
		db *sqlx.DB
	}

	Database interface {
		LoginUser(user *User, token *Token) (err error)
		GetToken(id uint64) (token *Token, err error)
		AddGuilds(guilds []*Guild) (err error)
		AddGuild(guild *Guild) (err error)
		AddUserGuilds(userGuilds []*UserGuild) (err error)
		GetBotExistsGuilds() (guilds []*Guild, err error)
		UpdateGuild(guild *Guild) (err error)
		UpdateGuildBotExists(id uint64, exists bool) (err error)
		GetBotAndUserExistsGuilds(userid uint64) ([]*Guild, error)
		AddEmojis(emojis []*discord.Emoji) (err error)
		AddEmoji(emoji *discord.Emoji) (err error)
		AddUserEmojis(userid uint64) (err error)
		GetUserEmojis(userid uint64) ([]*Emoji, error)
		GetEmojiUsers(userid uint64) ([]*EmojiUser, error)
		UpdateUserEmojis(obj UpdateEmojis) (err error)
		UpdateEmojiIfNotExists(guildid string, emojiids []string) (err error)
	}
)

func NewModel(db *sqlx.DB) *model {
	return &model{
		db: db,
	}
}

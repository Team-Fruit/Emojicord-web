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
		GetToken(id string) (token *Token, err error)
		AddGuilds(guilds *[]Guild) (err error)
		AddGuild(guild *Guild) (err error)
		AddUserGuilds(userGuilds *[]UserGuild) (err error)
		GetBotExistsGuilds() (guilds *[]Guild, err error)
		UpdateGuild(guild *Guild) (err error)
		UpdateGuildBotExists(id string, exists bool) (err error)
		AddEmojis(emojis *[]discord.Emoji) (err error)
		AddEmoji(emoji *discord.Emoji) (err error)
		AddUserEmojis(userEmojis *[]UserEmoji) (err error)
	}
)

func NewModel(db *sqlx.DB) *model {
	return &model{
		db: db,
	}
}

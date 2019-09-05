package model

import (
	"github.com/jmoiron/sqlx"
)


type (
	model struct {
		db *sqlx.DB
	}

	Database interface {
		LoginUser(user *User, token *Token) (err error)
		GetToken(id string) (token *Token, err error)
		AddGuilds(guilds *[]Guild) (err error)
		AddUserGuild(userGuilds *[]UserGuild) (err error)
		GetBotExistsGuilds() (guilds *[]Guild, err error)
		UpdateGuild(guild *Guild) (err error)
		UpdateGuildBotExists(id string, exists bool) (err error)
	}
)

func NewModel(db *sqlx.DB) *model {
	return &model{
		db: db,
	}
}

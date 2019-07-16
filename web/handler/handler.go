package handler

import (
	"github.com/Team-Fruit/Emojicord-web/web/model"
	"github.com/Team-Fruit/Emojicord-web/web/discord"
)

type handler struct {
	Model model.Database
	Bot   discord.BotClient
	User  discord.UserClient
}

func NewHandler(d model.Database, c discord.BotClient, u discord.UserClient) *handler {
	return &handler{
		Model: d,
		Bot: c,
		User: u,
	}
}
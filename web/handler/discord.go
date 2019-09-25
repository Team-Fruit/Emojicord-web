package handler

import (
	"fmt"
	
	"github.com/bwmarrin/discordgo"

	"github.com/Team-Fruit/Emojicord-web/web/model"
)

func (h *handler) GuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	fmt.Println("GuildCreate:", e.Guild.Name)

	guild := &model.Guild{
		ID: e.Guild.ID,
		Name: e.Guild.Name,
		Icon: e.Guild.Icon,
		BotExists: true,
	}
	
	if err := h.Model.AddGuild(guild); err != nil {
		fmt.Println("Failed to add guild", err)
	}

	emojis, err := h.Bot.GetEmojis(e.Guild.ID)
	if err != nil {
		fmt.Println("Failed to add get emoji", err)
	}

	if err := h.Model.AddEmojis(emojis); err != nil {
		fmt.Println("Failed to add emoji", err)
	}
}

func (h *handler) GuildUpdate(s *discordgo.Session, e *discordgo.GuildUpdate) {
	fmt.Println("GuildUpdate:", e.Guild.Name)

	guild := &model.Guild{
		ID: e.Guild.ID,
		Name: e.Guild.Name,	
		Icon: e.Guild.Icon,
		BotExists: true,
	}
	
	if err := h.Model.UpdateGuild(guild); err != nil {
		fmt.Println("Failed to update guild", err)
	}
}

func (h *handler) GuildDelete(s *discordgo.Session, e *discordgo.GuildDelete) {
	fmt.Println("GuildDelete:", e.Guild.Name)

	guilds, err := h.Bot.GetGuilds()
	if err != nil {
		fmt.Println("Failed to get bot guilds", err)
		return
	}

	exists := false
	for _, g := range *guilds {
		if g.ID == e.Guild.ID {
			exists = true
			break
		}
	}

	if !exists {
		if err := h.Model.UpdateGuildBotExists(e.Guild.ID, false); err != nil {
			fmt.Println("Failed to update guild", err)
		}
	}
}

func (h *handler) EmojisUpdate(s *discordgo.Session, e *discordgo.GuildEmojisUpdate) {
	fmt.Println("EmojisUpdate:", e.GuildID)
}
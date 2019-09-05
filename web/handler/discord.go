package handler

import (
	"fmt"
	
	"github.com/bwmarrin/discordgo"
)

func (h *handler) GuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	fmt.Println("GuildCreate:", e.Guild.Name)

	if err := h.Model.UpdateGuildBotExists(e.Guild.ID, true); err != nil {
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

}
package handler

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"

	"github.com/Team-Fruit/Emojicord-web/web/model"
)

func (h *handler) GuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	fmt.Println("GuildCreate:", e.Guild.Name)

	id, err := strconv.ParseUint(e.Guild.ID, 10, 64)
	if err != nil {
		fmt.Println("Failed to parse guild id", err)
		return
	}

	guild := &model.Guild{
		ID:        id,
		Name:      e.Guild.Name,
		Icon:      e.Guild.Icon,
		BotExists: true,
	}

	if err := h.Model.AddGuild(guild); err != nil {
		fmt.Println("Failed to add guild", err)
		return
	}

	emojis, err := h.Bot.GetEmojis(id)
	if err != nil {
		fmt.Println("Failed to add get emoji", err)
		return
	}

	if err := h.Model.AddEmojis(emojis); err != nil {
		fmt.Println("Failed to add emoji", err)
		return
	}
}

func (h *handler) GuildUpdate(s *discordgo.Session, e *discordgo.GuildUpdate) {
	fmt.Println("GuildUpdate:", e.Guild.Name)

	id, err := strconv.ParseUint(e.Guild.ID, 10, 64)
	if err != nil {
		fmt.Println("Failed to parse guild id", err)
		return
	}

	guild := &model.Guild{
		ID:        id,
		Name:      e.Guild.Name,
		Icon:      e.Guild.Icon,
		BotExists: true,
	}

	if err := h.Model.UpdateGuild(guild); err != nil {
		fmt.Println("Failed to update guild", err)
		return
	}
}

func (h *handler) GuildDelete(s *discordgo.Session, e *discordgo.GuildDelete) {
	fmt.Println("GuildDelete:", e.Guild.Name)

	id, err := strconv.ParseUint(e.Guild.ID, 10, 64)
	if err != nil {
		fmt.Println("Failed to parse guild id", err)
		return
	}

	if err := h.Model.UpdateGuildBotExists(id, false); err != nil {
		fmt.Println("Failed to update guild", err)
		return
	}
}

func (h *handler) EmojisUpdate(s *discordgo.Session, e *discordgo.GuildEmojisUpdate) {
	fmt.Println("EmojisUpdate:", e.GuildID)

	id, err := strconv.ParseUint(e.GuildID, 10, 64)
	if err != nil {
		fmt.Println("Failed to parse guild id", err)
		return
	}

		emojis, err := h.Bot.GetEmojis(id)
		if err != nil {
			fmt.Println("Failed to add get emoji", err)
			return
		}

		if err := h.Model.AddEmojis(emojis); err != nil {
			fmt.Println("Failed to add emoji", err)
			return
		}
}

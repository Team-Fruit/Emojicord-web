package handler

import (
	"fmt"
	
	"github.com/bwmarrin/discordgo"
)

func (h *handler) GuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	fmt.Println(e.Guild.Name)
}

func (h *handler) EmojisUpdate(s *discordgo.Session, e *discordgo.GuildEmojisUpdate) {

}
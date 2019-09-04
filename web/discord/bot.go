package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func InitBotUser(token string) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	discord.AddHandler(guildCreate)
	discord.AddHandler(emojisUpdate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
}

func guildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	fmt.Println(e.Guild.Name)
}

func emojisUpdate(s *discordgo.Session, e *discordgo.GuildEmojisUpdate) {

}

package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func LoginBotUser(token string) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	discord.AddHandler(emojisUpdate)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
}

func emojisUpdate(s* discordgo.Session, m *discordgo.GuildEmojisUpdate) {

}

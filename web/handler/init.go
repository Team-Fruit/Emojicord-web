package handler

import "fmt"

func (h *handler) Init() error {
	remoteGuilds, err := h.Bot.GetGuilds()
	if err != nil {
		return err
	}
	fmt.Println("SU")

	localGuilds, err := h.Model.GetBotExistsGuilds()
	if err != nil {
		return err
	}
	fmt.Println("SHI")

	for _, lg := range localGuilds {
		var exists = false
		for _, rg := range *remoteGuilds {
			if lg.ID == rg.ID {
				exists = true
				break
			}
		}

		if !exists {
			lg.BotExists = false
			if err := h.Model.UpdateGuild(&lg); err != nil {
				return err
			}
		}
	}
	return nil
}
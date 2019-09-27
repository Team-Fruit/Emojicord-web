package handler

import "fmt"

func (h *handler) Init() error {
	remoteGuilds, err := h.Bot.GetGuilds()
	if err != nil {
		return err
	}

	localGuilds, err := h.Model.GetBotExistsGuilds()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for i := range localGuilds {
		var exists = false
		for j := range remoteGuilds {
			if localGuilds[i].ID == remoteGuilds[j].ID {
				exists = true
				break
			}
		}

		if !exists {
			localGuilds[i].BotExists = false
			if err := h.Model.UpdateGuild(localGuilds[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
package handler

import (
	"github.com/Team-Fruit/Emojicord-web/web/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func (h *handler) GetEmojis(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	id := claims.ID

	token, err := h.Model.GetToken(id)
	if err != nil {
		return err
	}

	userGuilds, err := h.User.GetGuilds(token.ToOAuth2Token())
	if err != nil {
		return err
	}

	botGuilds, err := h.Model.GetBotExistsGuilds()
	if err != nil {
		return err
	}

	// Guilds Update
	// users__discord_guilds
	ug := make([]model.UserGuild, 0, len(*botGuilds))
	for _, g := range *userGuilds {
		permissions := uint(g.Permissions)
		canInvite := (permissions & 0x20) == 0x20

		for _, b := range *botGuilds {
			if g.ID == b.ID {
				ug = append(ug, model.UserGuild{
					UserID:      id,
					GuildID:     g.ID,
					IsOwner:     g.Owner,
					Permissions: permissions,
					CanInvite:   canInvite,
				})
				break
			}
		}
	}

	err = h.Model.AddUserGuilds(&ug)
	if err != nil {
		return err
	}

	return c.String(200, "OK")
}

package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/Team-Fruit/Emojicord-web/web/model"
)

type Guild struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	Owner     bool   `json:"owner"`
	BotExists bool   `json:"botexists"`
	CanInvite bool   `json:"caninvite"`
}

func (h *handler) GetGuilds(c echo.Context) error {
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

	botGuilds, err := h.Bot.GetGuilds()
	if err != nil {
		return err
	}

	nug := make([]model.Guild, 0, len(*userGuilds))
	ug := make([]model.UserGuild, 0, len(*userGuilds))
	rg := make([]Guild, 0, len(*userGuilds))
	for _, g := range *userGuilds {
		permissions := uint(g.Permissions)
		canInvite := (permissions & 0x20) == 0x20

		nug = append(nug, model.Guild{
			ID:        g.ID,
			Name:      g.Name,
			Icon:      g.Icon,
			BotExists: false,
		})

		ug = append(ug, model.UserGuild{
			UserID:      id,
			GuildID:     g.ID,
			IsOwner:     g.Owner,
			Permissions: permissions,
			CanInvite:   canInvite,
		})

		rg = append(rg, Guild{
			ID:        g.ID,
			Name:      g.Name,
			Icon:      g.Icon,
			Owner:     g.Owner,
			CanInvite: canInvite,
		})
	}

	err = h.Model.AddGuilds(&nug)
	if err != nil {
		return err
	}

	err = h.Model.AddUserGuild(&ug)
	if err != nil {
		return err
	}

	//Update Bot Guilds
	nbg := make([]model.Guild, 0, len(*botGuilds))
	for _, g := range *botGuilds {
		nbg = append(nbg, model.Guild{
			ID:        g.ID,
			Name:      g.Name,
			Icon:      g.Icon,
			BotExists: true,
		})

		for i := range rg {
			if rg[i].ID == g.ID {
				rg[i].BotExists = true
				break
			}
		}
	}

	err = h.Model.AddGuilds(&nbg)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rg)
}

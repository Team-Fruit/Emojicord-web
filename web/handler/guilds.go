package handler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/Team-Fruit/Emojicord-web/web/model"
)

type Guild struct {
	ID        uint64 `json:"id,string"`
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

	botGuilds, err := h.Model.GetBotExistsGuilds()
	if err != nil {
		return err
	}

	// discord_guilds
	nug := make([]*model.Guild, 0, len(userGuilds))
	// users__discord_guilds
	ug := make([]*model.UserGuild, 0, len(userGuilds))
	// Response
	rg := make([]*Guild, 0, len(userGuilds))
	for i := range userGuilds {
		g := userGuilds[i]

		permissions := uint(g.Permissions)
		canInvite := (permissions & 0x20) == 0x20

		botExists := false
		for j := range botGuilds {
			b := botGuilds[j]

			if g.ID == b.ID {
				botExists = true
				break
			}
		}

		if !botExists {
			nug = append(nug, &model.Guild{
				ID:        g.ID,
				Name:      g.Name,
				Icon:      g.Icon,
				BotExists: false,
			})
		}

		ug = append(ug, &model.UserGuild{
			UserID:      id,
			GuildID:     g.ID,
			IsOwner:     g.Owner,
			Permissions: permissions,
			CanInvite:   canInvite,
		})

		rg = append(rg, &Guild{
			ID:        g.ID,
			Name:      g.Name,
			Icon:      g.Icon,
			Owner:     g.Owner,
			BotExists: botExists,
			CanInvite: canInvite,
		})
	}

	err = h.Model.AddGuilds(nug)
	if err != nil {
		return err
	}

	err = h.Model.AddUserGuilds(ug)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, rg)
}

package handler

import (
	"net/http"

	"github.com/Team-Fruit/Emojicord-web/web/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	Emojis struct {
		Emoji     []*model.Emoji     `json:"emojis"`
		Guild     []*model.Guild     `json:"guilds"`
		EmojiUser []*model.EmojiUser `json:"users"`
	}
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
	ug := make([]*model.UserGuild, 0, len(botGuilds))
	for i := range userGuilds {
		g := userGuilds[i]

		permissions := uint(g.Permissions)
		canInvite := (permissions & 0x20) == 0x20

		for j := range botGuilds {
			b := botGuilds[j]

			if g.ID == b.ID {
				ug = append(ug, &model.UserGuild{
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

	err = h.Model.AddUserGuilds(ug)
	if err != nil {
		return err
	}

	res := Emojis{}

	guild, err := h.Model.GetBotAndUserExistsGuilds(id)
	if err != nil {
		return err
	}

	res.Guild = guild

	if err := h.Model.AddUserEmojis(id); err != nil {
		return err
	}

	emoji, err := h.Model.GetUserEmojis(id)
	if err != nil {
		return err
	}

	res.Emoji = emoji

	emojiUser, err := h.Model.GetEmojiUsers(id)
	if err != nil {
		return err
	}

	res.EmojiUser = emojiUser

	return c.JSON(http.StatusOK, res)
}

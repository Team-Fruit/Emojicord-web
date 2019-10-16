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

func (h *handler) PutEmoji(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	id := claims.ID

	emojiid := c.Param("id")

	if err := h.Model.UpdateUserEmoji(model.UpdateEmoji{
		Enabled: true,
		UserID:  id,
		EmojiID: emojiid,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) DeleteEmoji(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	id := claims.ID

	emojiid := c.Param("id")

	if err := h.Model.UpdateUserEmoji(model.UpdateEmoji{
		Enabled: false,
		UserID:  id,
		EmojiID: emojiid,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) PutEmojis(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	id := claims.ID

	emojiids := []string{}
	if err := c.Bind(&emojiids); err != nil {
		return err
	}

	if err := h.Model.UpdateUserEmojis(model.UpdateEmojis{
		Enabled: true,
		UserID:  id,
		EmojiID: emojiids,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *handler) DeleteEmojis(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	id := claims.ID

	emojiids := []string{}
	if err := c.Bind(&emojiids); err != nil {
		return err
	}

	if err := h.Model.UpdateUserEmojis(model.UpdateEmojis{
		Enabled: false,
		UserID:  id,
		EmojiID: emojiids,
	}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

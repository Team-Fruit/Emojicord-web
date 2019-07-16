package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
)

func (h *handler) GetGuilds(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	id := claims.ID

	token, err := h.Model.GetToken(id)
	if err != nil {
		return err
	}

	

	guilds, err := h.Bot.GetGuilds()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, guilds)
}


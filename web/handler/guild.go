package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
)

func (h *handler) GetGuilds(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	name := claims.Username
	return c.String(http.StatusOK, "Welcome back "+name+"!")
}

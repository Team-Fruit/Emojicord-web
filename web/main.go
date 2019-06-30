package main

import (
	"os"
	"net/http"
	"time"
	"math/rand"
	"context"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/oauth2"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", index)
	e.GET("/callback", callback)

	e.Logger.Fatal(e.Start(":8082"))
}

const cookieName = "discordOAuth2State"

type Guild struct {
	Owner       bool   `json:"owner"`
	Permissions int    `json:"permissions"`
	Icon        string `json:"icon"`
	ID          string `json:"id"`
	Name        string `json:"name"`
}

func GetConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint:     oauth2.Endpoint{
			AuthURL: "https://discordapp.com/api/oauth2/authorize",
			TokenURL: "https://discordapp.com/api/oauth2/token",
		},
		Scopes:       []string{"identify", "guilds"},
		RedirectURL:  "https://emojicord.teamfruit.net/api/callback",
	}
}

func index(c echo.Context) error {
	config := GetConfig()

	state, err := generateState()
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		Secure:   true,
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	url := config.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "code"))
	return c.Redirect(http.StatusFound, url)
}

func callback(c echo.Context) error {
	cookieState, err := c.Cookie(cookieName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to read Cookie")
	}
	paramState := c.QueryParam("state")
	if cookieState.Value != paramState {
		return echo.NewHTTPError(http.StatusBadRequest, "OAuth2 State Error")
	}

	if e := c.QueryParam("error"); e != "" {
		return echo.NewHTTPError(http.StatusUnauthorized, c.QueryParam("error_description"))
	}

	config := GetConfig()
	token, err := config.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		return err
	}

	client := config.Client(context.Background(), token)
	res, err := client.Get("https://discordapp.com/api/v6/users/@me/guilds")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	
	var guilds []Guild
	if err := json.Unmarshal(body, &guilds); err != nil {
		return err;
	}
	return c.JSON(http.StatusOK, guilds)
}

func generateState() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

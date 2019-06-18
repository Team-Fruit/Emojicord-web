package main

import (
	"fmt"
	"os"
	"math/rand"
	"context"
	"io/ioutil"
	"encoding/json"

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
		Scopes:       []string{"guilds"},
		RedirectURL:  "https://emojicord.teamfruit.net/callback",
	}
}

func index(c echo.Context) error {
	config := GetConfig()

	stateBytes := make([]byte, 16)
	_, err := rand.Read(stateBytes)
	if err != nil {
		return err
	}
	state := fmt.Sprintf("%x", stateBytes)
	url := config.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "code"))
	return c.Redirect(302, url)
}

func callback(c echo.Context) error {
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
	
	fmt.Println(string(body))
	var guilds []Guild
	if err := json.Unmarshal(body, &guilds); err != nil {
		return err;
	}
	return c.JSON(200, guilds)
}

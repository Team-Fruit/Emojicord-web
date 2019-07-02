package handler

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
	"golang.org/x/oauth2"
	"github.com/dgrijalva/jwt-go"
)

const cookieName = "discordOAuth2State"

// type Guild struct {
// 	Owner       bool   `json:"owner"`
// 	Permissions int    `json:"permissions"`
// 	Icon        string `json:"icon"`
// 	ID          string `json:"id"`
// 	Name        string `json:"name"`
// }

type JWTClaims struct {
	Username      string `json:"username"`
	Locale        string `json:"locale"`
	Avater        string `json:"avater"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
	jwt.StandardClaims
}

type User struct {
	Username      string `json:"username"`
	Locale        string `json:"locale"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Flags         int    `json:"flags"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
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

func Auth(c echo.Context) error {
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
	return c.Redirect(http.StatusSeeOther, url)
}

func Callback(c echo.Context) error {
	cookieState, err := c.Cookie(cookieName)
	if err != nil {
		return err
	}
	paramState := c.QueryParam("state")
	if cookieState.Value != paramState {
		return err
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
	res, err := client.Get("https://discordapp.com/api/v6/users/@me")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	
	// var guilds []Guild
	// if err := json.Unmarshal(body, &guilds); err != nil {
	// 	return err;
	// }
	// return c.JSON(http.StatusOK, guilds)
	
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return err;
	}

	claims := &JWTClaims{
		user.Username,
		user.Locale,
		user.Avatar,
		user.Discriminator,
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, s)
}

func generateState() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

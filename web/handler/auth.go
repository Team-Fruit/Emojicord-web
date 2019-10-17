package handler

import (
	"os"
	"net/http"
	"net/url"
	"time"
	"math/rand"
	"context"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"github.com/dgrijalva/jwt-go"

	"github.com/Team-Fruit/Emojicord-web/web/model"
)

const cookieName = "discordOAuth2State"

type JWTClaims struct {
	Username      string `json:"username"`
	Locale        string `json:"locale"`
	Avater        string `json:"avater"`
	Discriminator string `json:"discriminator"`
	ID            uint64 `json:"id,string"`
	jwt.StandardClaims
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
		RedirectURL:  os.Getenv("REDIRECT_URL"),
	}
}

func (h *handler) Auth(c echo.Context) error {
	config := GetConfig()

	state, err := generateState()
	if err != nil {
		return c.Redirect(http.StatusSeeOther, createErrorRedirectURL("internal_server_error", "Internal Server Error"))
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

func (h *handler) Callback(c echo.Context) error {
	cookieState, err := c.Cookie(cookieName)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, createErrorRedirectURL("internal_server_error", "Internal Server Error"))
	}
	paramState := c.QueryParam("state")
	if cookieState.Value != paramState {
		return c.Redirect(http.StatusSeeOther, createErrorRedirectURL("internal_server_error", "Internal Server Error"))
	}

	if e := c.QueryParam("error"); e != "" {
		return c.Redirect(http.StatusSeeOther, createErrorRedirectURL(c.QueryParam("error"), c.QueryParam("error_description")))
	}

	config := GetConfig()
	token, err := config.Exchange(context.Background(), c.QueryParam("code"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, createErrorRedirectURL("internal_server_error", "Internal Server Error"))
	}

	user, err := getUserFromDiscord(config, token)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, createErrorRedirectURL("internal_server_error", "Internal Server Error"))
	}
	
	mt := model.ToModelToken(user.ID, token)

	if err := h.Model.LoginUser(user, mt); err != nil {
		return c.Redirect(http.StatusSeeOther, createErrorRedirectURL("internal_server_error", "Internal Server Error"))
	}

	t, err := generateJWT(user)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, createErrorRedirectURL("internal_server_error", "Internal Server Error"))
	}

	return c.Redirect(http.StatusSeeOther, os.Getenv("LOGIN_URL") + "?token=" + t)
}

func generateState() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func createErrorRedirectURL(err string, desc string) string {
	params := url.Values{}
	params.Add("error", err)
	params.Add("error_description", desc)
	return os.Getenv("LOGIN_URL") + "?" + params.Encode()
}

func generateJWT(user *model.User) (token string, err error) {
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

	return t.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func getUserFromDiscord(config *oauth2.Config, token *oauth2.Token) (user *model.User, err error) {
	client := config.Client(context.Background(), token)
	
	res, err := client.Get("https://discordapp.com/api/v6/users/@me")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &user); err != nil {
		return nil, err
	}

	return user, nil
}
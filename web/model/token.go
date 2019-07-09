package model

import (
	"time"
	"golang.org/x/oauth2"
)

type Token struct {
	UserID       string    `db:"user_id"`
	AccessToken  string    `db:"access_token"`
	TokenType    string    `db:"token_type"`
	RefreshToken string    `db:"refresh_token"`
	Expiry       time.Time `db:"expiry"`
}

func (t *Token) ToOAuth2Token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: t.AccessToken,
		TokenType: t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry: t.Expiry,
	}
}

func ToModelToken(id string, t *oauth2.Token) *Token {
	return &Token{
		UserID: id,
		AccessToken: t.AccessToken,
		TokenType: t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry: t.Expiry,
	}
}
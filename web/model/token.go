package model

import (
	"time"
	"golang.org/x/oauth2"
)

type Token struct {
	UserID       uint64    `db:"user_id"`
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

func ToModelToken(id uint64, t *oauth2.Token) *Token {
	return &Token{
		UserID: id,
		AccessToken: t.AccessToken,
		TokenType: t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry: t.Expiry,
	}
}

func (m *model) GetToken(id uint64) (*Token, error) {
	var token Token
	if err := m.db.Get(&token, "SELECT * FROM users__discord_tokens WHERE user_id = ?", id); err != nil {
		return nil, err
	}
	return &token, nil
}
package model

type User struct {
		ID            string `json:"id" db:"id"`
		Username      string `json:"username" db:"username"`
		Discriminator string `json:"discriminator" db:"discriminator"`
		Avatar        string `json:"avatar" db:"avatar"`
		Locale        string `json:"locale" db:"locale"`
		CreatedAt     string `db:"created_at"`
		LastLogin     string `db:"last_login"`
}

func (m *model) LoginUser(user *User, token *Token) (err error) {
	tx, err := m.db.Begin()
    if err != nil {
        return
    }
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
        err = tx.Commit()
	}()
	
	_, err = tx.Exec(`INSERT INTO users
					(id, username, discriminator, avatar, locale)
					VALUES (?, ?, ?, ?, ?)
					ON DUPLICATE KEY UPDATE
					username = VALUES(username),
					discriminator = VALUES(discriminator),
					avatar = VALUES(avatar),
					locale = VALUES(locale),
					last_login = NOW()`,
					user.ID,
					user.Username,
					user.Discriminator,
					user.Avatar,
					user.Locale)
	if err != nil {
		return
	}

	_, err = tx.Exec(`INSERT INTO users__discord_tokens
					(user_id, access_token, token_type, refresh_token, expiry)
					VALUES (?, ?, ?, ?, ?)
					ON DUPLICATE KEY UPDATE
					access_token = VALUES(access_token),
					token_type = VALUES(token_type),
					refresh_token = VALUES(refresh_token),
					expiry = VALUES(expiry)`,
					token.UserID,
					token.AccessToken,
					token.TokenType,
					token.RefreshToken,
					token.Expiry)
	if err != nil {
		return
	}

	return
}
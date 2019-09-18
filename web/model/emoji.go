package model

type (
	Emoji struct {
		ID       string `db:"id"`
		GuildID  string `db:"guild_id"`
		UserID   string `db:"user_id"`
		Name     string `db:"name"`
		Animated bool   `db:"is_animated"`
	}

	UserEmoji struct {
		UserID  string `db:"user_id"`
		EmojiID string `db:"emoji_id"`
		Enabled bool   `db:"is_enabled"`
	}
)

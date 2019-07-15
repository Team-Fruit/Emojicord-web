package model

type Guild struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Icon      string `json:"icon" db:"icon"`
	BotExists bool   `db:"is_bot_exists"`
}
package model

type User struct {
	Username      string `json:"username"`
	Locale        string `json:"locale"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
}
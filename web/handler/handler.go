package handler

import (
	"github.com/Team-Fruit/Emojicord-web/web/model"
)

type handler struct {
	Model model.Database
}

func NewHandler(d models.Database) *handler {
	return &handler{
		Model: d,
	}
}
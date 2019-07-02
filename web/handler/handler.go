package handler

import (
	"github.com/Team-Fruit/Emojicord-web/web/model"
)

type handler struct {
	Model model.Database
}

func NewHandler(d model.Database) *handler {
	return &handler{
		Model: d,
	}
}
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/Team-Fruit/Emojicord-web/web/handler"
	"github.com/Team-Fruit/Emojicord-web/web/model"
)

func main() {
	db := sqlx.MustConnect("mysql", "emojicord:@tcp(db:3306)/emojicord_db")
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	h := handler.NewHandler(model.NewModel(db))

	e.GET("/auth/login", h.Auth)
	e.GET("/auth/callback", h.Callback)

	e.Logger.Fatal(e.Start(":8082"))
}

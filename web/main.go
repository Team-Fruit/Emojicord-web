package main

import (
	"github.com/jmoiron/sqlx"
	
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/Team-Fruit/Emojicord-web/web/handler"
)

func main() {
	db := sqlx.MustConnect("mysql", "emojicord:@tcp(db:3306)/emojicord_db")
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/login", handler.Auth)
	e.GET("/callback", handler.Callback)

	e.Logger.Fatal(e.Start(":8082"))
}

package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/Team-Fruit/Emojicord-web/web/handler"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/login", handler.Auth)
	e.GET("/callback", handler.Callback)

	e.Logger.Fatal(e.Start(":8082"))
}

package main

import (
	"os"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/bwmarrin/discordgo"

	"github.com/Team-Fruit/Emojicord-web/web/handler"
	"github.com/Team-Fruit/Emojicord-web/web/model"
	"github.com/Team-Fruit/Emojicord-web/web/discord"
)

func main() {
	db := sqlx.MustConnect("mysql", "emojicord:password@tcp(db:3306)/emojicord_db?parseTime=true")
	defer db.Close()

	m := model.NewModel(db)
	b := discord.NewBotClient(os.Getenv("BOT_TOKEN"))
	u := discord.NewUserClient(handler.GetConfig())
	h := handler.NewHandler(m, b, u)

	dc, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("Failed to creating Discord session,", err)
	}

	dc.AddHandler(h.GuildCreate)
	dc.AddHandler(h.EmojisUpdate)

	err = dc.Open()
	if err != nil {
		fmt.Println("Failed to opening discord bot connection,", err)
		return
	}


	if err := h.Init(); err != nil {
		fmt.Println("Failed to guild update initialization", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/auth/login", h.Auth)
	e.GET("/auth/callback", h.Callback)

	g := e.Group("/user")
	config := middleware.JWTConfig{
		Claims: &handler.JWTClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	g.Use(middleware.JWTWithConfig(config))
	g.GET("/guilds", h.GetGuilds)

	e.Logger.Fatal(e.Start(":8082"))
}

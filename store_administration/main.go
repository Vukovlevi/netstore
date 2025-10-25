package main

import (
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/vukovlevi/netstore/store_administration/config"
	"github.com/vukovlevi/netstore/store_administration/db"
	"github.com/vukovlevi/netstore/store_administration/route"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        slog.Error("could not load environment variables", "error", err)
        panic("unable to load environment variables")
    }

    config := config.CreateApplicationConfig()
    err = config.Apply()
    if err != nil {
        slog.Error("unable to apply application config", "error", err)
        panic("unable to apply application config")
    }
    defer db.Disconnect()

    e := echo.New()
    apiGroup := e.Group("/api")
    apiGroup.POST("/login", route.HandleLogin)
    e.GET("/", func(c echo.Context) error {return c.JSON(http.StatusOK, map[string]string{"message": "itt vagy"})})

    e.Logger.Fatal(e.Start(":8000")) //TODO: read address from config
}

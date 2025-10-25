package main

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/vukovlevi/netstore/store_administration/config"
	"github.com/vukovlevi/netstore/store_administration/db"
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
}

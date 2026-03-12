package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/config"
	"github.com/vukovlevi/netstore/store_administration/db"
	"github.com/vukovlevi/netstore/store_administration/middleware"
	"github.com/vukovlevi/netstore/store_administration/network"
	"github.com/vukovlevi/netstore/store_administration/route"

	mw "github.com/labstack/echo/v4/middleware"
)

func main() {
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

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

    err = network.NewNetworkManager(config.CentralServerAddress, config.CentralServerPort, config.Psk)
    if err != nil {
        slog.Error("could not create network manager", "error", err)
    }

    e := echo.New()
    e.Use(mw.CORSWithConfig(mw.CORSConfig{
        AllowOrigins: []string{"http://localhost"}, //TODO: read from config
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    apiGroup := e.Group("/api")
    apiAuthGroup := apiGroup.Group("", middleware.AuthenticateUser)
    apiStoreLeaderGroup := apiAuthGroup.Group("", middleware.AuthorizeStoreLeader)
    apiStoreLeaderOrHRGroup := apiAuthGroup.Group("", middleware.AuthorizeStoreLeaderOrHR)

    apiGroup.POST("/login", route.HandleLogin)
    apiAuthGroup.GET("/logout", route.HandleLogout)

    apiStoreLeaderOrHRGroup.GET("/contract-type", route.HandleGetAllContractType)
    apiStoreLeaderGroup.POST("/contract-type", route.HandlePostContractType)
    apiStoreLeaderGroup.PUT("/contract-type", route.HandleUpdateContractType)
    apiStoreLeaderGroup.DELETE("/contract-type", route.HandleDeleteContractType)

    apiStoreLeaderGroup.GET("/store-type", route.HandleGetAllStoreType)
    apiStoreLeaderGroup.GET("/store-detail", route.HandleGetStoreDetail)
    apiStoreLeaderGroup.PUT("/store-detail", route.HandleUpdateStoreDetail)

    apiAuthGroup.GET("/user", route.HandleGetUser)
    apiStoreLeaderOrHRGroup.GET("/all-user", route.HandleGetAllUser)
    apiStoreLeaderOrHRGroup.POST("/user", route.HandlePostUser)
    apiStoreLeaderOrHRGroup.PUT("/user", route.HandleUpdateUser)
    apiStoreLeaderOrHRGroup.DELETE("/user", route.HandleDeleteUser)
    apiStoreLeaderOrHRGroup.GET("/role", route.HandleGetAllRole)
    apiAuthGroup.POST("/password-change", route.HandleUpdateUserPassword)

    apiStoreLeaderGroup.GET("/open-hour", route.HandleGetOpenHours)
    apiStoreLeaderGroup.POST("/open-hour", route.HandlePostOpenHour)
    apiStoreLeaderGroup.PUT("/open-hour", route.HandleUpdateOpenHour)
    apiStoreLeaderGroup.DELETE("/open-hour", route.HandleDeleteOpenHour)
    apiStoreLeaderGroup.GET("/weekdays", route.HandleGetWeekDays)

    apiStoreLeaderGroup.GET("/connect", route.HandleGetConnect)
    apiStoreLeaderGroup.POST("/connect", route.HandlePostConnect)
    apiAuthGroup.POST("/network-search", route.HandlePostNetworkSearch)

    apiStoreLeaderOrHRGroup.GET("/contract", route.HandleGetContractByUserId)
    apiStoreLeaderOrHRGroup.POST("/contract", route.HandlePostContract)
    apiStoreLeaderOrHRGroup.PUT("/contract", route.HandleUpdateContract)
    apiStoreLeaderOrHRGroup.DELETE("/contract", route.HandleDeleteContract)
    apiStoreLeaderOrHRGroup.GET("/contract-file", route.HandleGetContractFile)
    apiStoreLeaderOrHRGroup.DELETE("/contract-file", route.HandleDeleteContractFile)

    apiAuthGroup.GET("/echo", route.HandleGetEcho)
    e.Static("/assets", "public/assets")
    e.GET("/*", func(c echo.Context) error {return c.File("public/index.html")})

    apiAuthGroup.GET("/", func(c echo.Context) error {return c.JSON(http.StatusOK, map[string]string{"message": "itt vagy"})})

    e.Logger.Fatal(e.Start(config.ToAddress()))
}

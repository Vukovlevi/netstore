package route

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/network"
)

type ConnectInfo struct {
	Ip string `json:"ip_address"`
	Port string `json:"port"`
	Psk string `json:"psk"`
}

func HandleGetConnect(c echo.Context) error {
	if network.Manager.IsConnected() {
		return c.JSON(http.StatusOK, CreateMessage("A rendszer csatlakoztatva van a központi szerverhez!"))
	} else {
		return c.JSON(http.StatusOK, CreateMessage("A rendszer nincs csatlakoztatva a központi szerverhez!"))
	}
}

func HandlePostConnect(c echo.Context) error {
	connectInfo := ConnectInfo{}
	if err := c.Bind(&connectInfo); err != nil {
		slog.Error("could not bind connect info", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("Csatlakozási adatok olvasása sikertelen!"))
	}

	if network.Manager.IsConnected() {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage("A rendszer már csatlakozva van a központi szerverhez!"))
	}

	err := network.NewNetworkManager(connectInfo.Ip, connectInfo.Port, connectInfo.Psk)
	if err != nil {
		slog.Error("could not connect to central server", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, CreateMessage("Sikeres csatlakozás a központi szerverhez!"))
}

func HandlePostNetworkSearch(c echo.Context) error {
	if !network.Manager.IsConnected() {
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A rendszer nincs csatlakoztatva a központi szerverhez!"))
	}

	searchParam, err := io.ReadAll(c.Request().Body)
	c.Request().Body.Close()
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage("A hálózati keresés paraméterét nem sikerült kiolvasni!"))
	}

	searchResult, err := network.Manager.SearchNetwork(searchParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage(err.Error()))
	}

	return c.JSONBlob(http.StatusOK, searchResult)
}
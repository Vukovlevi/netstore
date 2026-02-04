package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandleGetAllRole(c echo.Context) error {
    roles, err := model.GetAllRole()
    if err != nil {
        slog.Error("coul not get all role", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A rangok lekérdezése sikertelen!"))
    }

    return c.JSON(http.StatusOK, roles)
}

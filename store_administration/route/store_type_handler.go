package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandleGetAllStoreType(c echo.Context) error {
    storeTypes, err := model.GetAllStoreType()
    if err != nil {
        slog.Error("could not get all store types", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("Az üzlet típusok lekérdezése sikertelen!"))
    }
    return c.JSON(http.StatusOK, storeTypes)
}

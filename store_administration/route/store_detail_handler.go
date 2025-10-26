package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandleGetStoreDetail(c echo.Context) error {
    storeDetail, err := model.GetStoreDetail()
    if err != nil {
        slog.Error("could not get store detail", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not get store detail")) //TODO: user-readable error message
    }
    return c.JSON(http.StatusOK, storeDetail)
}

func HandleUpdateStoreDetail(c echo.Context) error {
    storeDetail := model.StoreDetail{}
    if err := c.Bind(&storeDetail); err != nil {
        slog.Error("could not bind store detail", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not bind store detail")) //TODO: user-readable error message
    }

    if err := storeDetail.UpdateStoreDetail(); err != nil {
        slog.Error("could not update store detail", "error", err, "store detail", storeDetail)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("could not update store detail")) //TODO: user-readable error message
    }

    return c.JSON(http.StatusOK, CreateMessage("store detail updated")) //TODO: user-readable message
}

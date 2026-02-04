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
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("Az üzlet adatainak lekérdezése sikertelen!"))
    }
    return c.JSON(http.StatusOK, storeDetail)
}

func HandleUpdateStoreDetail(c echo.Context) error {
    storeDetail := model.StoreDetail{}
    if err := c.Bind(&storeDetail); err != nil {
        slog.Error("could not bind store detail", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("Az üzlet adatainak olvasása sikertelen!"))
    }

    if err := storeDetail.ValidateUpdate(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    if err := storeDetail.UpdateStoreDetail(); err != nil {
        slog.Error("could not update store detail", "error", err, "store detail", storeDetail)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("Az üzlet adatainak módosítása sikertelen!"))
    }

    return c.JSON(http.StatusOK, CreateMessage("Az üzlet adatainak módosítása sikeres!"))
}

package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandlePostContractType(c echo.Context) error {
    contractType := model.ContractType{}
    if err := c.Bind(&contractType); err != nil {
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződéstípus adatainak olvasása nem sikerült!"))
    }

    if err := contractType.ValidateInsert(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    if err := contractType.InsertNewContractType(); err != nil {
        slog.Error("could not save new contract type", "error", err, "contract type", contractType)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("Az új szerződéstípus mentése nem sikerült (lehet hogy már létezik ilyen néven)!"))
    }

    return c.JSON(http.StatusCreated, CreateMessage("Új szerződéstípus sikeresen létrehozva!"))
}

func HandleGetAllContractType(c echo.Context) error {
    contractTypes, err := model.GetAllContractType()
    if err != nil {
        slog.Error("could not get all contract types", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződéstípusok lekérdezése nem sikerült!"))
    }
    return c.JSON(http.StatusOK, contractTypes)
}

func HandleUpdateContractType(c echo.Context) error {
    contractType := model.ContractType{}
    if err := c.Bind(&contractType); err != nil {
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződéstípus adatainak olvasása nem sikerült!"))
    }

    if err := contractType.ValidateUpdate(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    if err := contractType.UpdateContractType(); err != nil {
        slog.Error("could not update contract type", "error", err, "contract type", contractType)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("A szerződéstípus módosítása sikertelen (lehet hogy már létezik típus ilyen néven)!"))
    }

    return c.JSON(http.StatusOK, CreateMessage("A szerződéstípus módosítása sikeres!"))
}

func HandleDeleteContractType(c echo.Context) error {
    contractType := model.ContractType{}
    if err := c.Bind(&contractType); err != nil {
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződéstípus adatainak olvasása sikertelen!"))
    }

    if err := contractType.ValidateDelete(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    if err := contractType.DeleteContractType(); err != nil {
        slog.Error("could not delete contract type", "error", err, "contract type", contractType)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("A szerződéstípus törlése sikertelen!"))
    }

    return c.NoContent(http.StatusNoContent)
}

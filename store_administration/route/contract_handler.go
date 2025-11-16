package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandleGetContracts(c echo.Context) error {
	contracts, err := model.GetAllContract()
	if err != nil {
		slog.Error("could not get contracts", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződések lekérdezése nem sikerült!"))
	}

	return c.JSON(http.StatusOK, contracts)
}

func HandlePostContract(c echo.Context) error {
	contract := model.Contract{}
	if err := c.Bind(&contract); err != nil {
		slog.Error("could not bind contract", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés adatait nem sikerült értelmezni!"))
	}

	if err := contract.ValidateInsert(); err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
	}

	if err := contract.InsertNewContract(); err != nil {
		slog.Error("could not insert new contract", "error", err, "contract", contract)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés mentése nem sikerült!"))
	}

	return c.JSON(http.StatusCreated, CreateMessage("A szerződés mentése sikeres!"))
}

func HandleUpdateContract(c echo.Context) error {
	contract := model.Contract{}
	if err := c.Bind(&contract); err != nil {
		slog.Error("could not bind contract", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés adatait nem sikerült értelmezni!"))
	}

	if err := contract.ValidateUpdate(); err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
	}

	if err := contract.UpdateContract(); err != nil {
		slog.Error("could not update contract", "error", err, "contract", contract)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés módosítása nem sikerült!"))
	}

	return c.JSON(http.StatusCreated, CreateMessage("A szerződés módosítása sikeres!"))
}

func HandleDeleteContract(c echo.Context) error {
	contract := model.Contract{}
	if err := c.Bind(&contract); err != nil {
		slog.Error("could not bind contract", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés adatait nem sikerült értelmezni!"))
	}

	if err := contract.ValidateDelete(); err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
	}

	if err := contract.DeleteContract(); err != nil {
		slog.Error("could not delete contract", "error", err, "contract", contract)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés törlése nem sikerült!"))
	}

	return c.NoContent(http.StatusNoContent)
}
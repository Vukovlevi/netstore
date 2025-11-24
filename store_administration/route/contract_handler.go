package route

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandleGetContractByUserId(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("userId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage("Érvénytelen felhasználó azonosító!"))
	}
	contract, err := model.GetContractByUserId(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		slog.Error("could not get contract for userId", "error", err, "userId", userId)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés lekérdezése nem sikerült!"))
	}

	return c.JSON(http.StatusOK, contract)
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

	return c.JSON(http.StatusOK, CreateMessage("A szerződés módosítása sikeres!"))
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
package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandleGetOpenHours(c echo.Context) error {
	openHours, err := model.GetOpenHours(false)
	if err != nil {
		slog.Error("could not get open hours", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A nyitvatartási idők lekérdezése sikertelen!"))
	}

	return c.JSON(http.StatusOK, openHours)
}

func HandlePostOpenHour(c echo.Context) error {
	openHour := model.OpenHour{}
	if err := c.Bind(&openHour); err != nil {
		slog.Error("could not bind open hour", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("Nem sikerült értelmezni a nyitvatartási időt!"))
	}

	if err := openHour.ValidateInsert(); err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
	}

	if err := openHour.InsertNewOpenHour(); err != nil {
		slog.Error("could not insert new open hour", "error", err, "open hour", openHour)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("Nem sikerült menteni a nyitvatartási időt!"))
	}

	return c.JSON(http.StatusCreated, CreateMessage("A nyitvatartási idő mentése sikeres!"))
}

func HandleUpdateOpenHour(c echo.Context) error {
	openHour := model.OpenHour{}
	if err := c.Bind(&openHour); err != nil {
		slog.Error("could not bind open hour", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A nyitvatartási időt nem sikerült értelmezni!"))
	}

	if err := openHour.ValidateUpdate(); err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
	}

	if err := openHour.UpdateOpenHour(); err != nil {
		slog.Error("could not update open hour", "error", err, "open hour", openHour)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A nyitvatartási időt nem sikerült frissíteni!"))
	}

	return c.JSON(http.StatusOK, CreateMessage("A nyitvatartási idő frissítése sikeres!"))
}

func HandleDeleteOpenHour(c echo.Context) error {
	openHour := model.OpenHour{}
	if err := c.Bind(&openHour); err != nil {
		slog.Error("could not bind open hour", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A nyitvatartási idő adatait nem sikerült értelmezni!"))
	}

	if err := openHour.ValidateDelete(); err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
	}

	if err := openHour.DeleteOpenHour(); err != nil {
		slog.Error("could not delete open hour", "error", err, "open hour", openHour)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A nyitvatartási idő törlése nem sikerült!"))
	}

	return c.NoContent(http.StatusNoContent)
}
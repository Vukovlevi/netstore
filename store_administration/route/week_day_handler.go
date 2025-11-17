package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandleGetWeekDays(c echo.Context) error {
	weekDays, err := model.GetAllWeekDay()
	if err != nil {
		slog.Error("could not get weekdays", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A hét napjait nem sikerült lekérdezni!"))
	}
	return c.JSON(http.StatusOK, weekDays)
}
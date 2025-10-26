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
        return c.JSON(http.StatusInternalServerError, CreateErrormessage("could not bind contract type")) //TODO: user-readalbe error message
    }

    if err := contractType.InsertNewContractType(); err != nil {
        slog.Error("could not save new contract type", "error", err, "contract type", contractType)
        return c.JSON(http.StatusBadRequest, CreateErrormessage("could not save contract type (possibly existing)")) //TODO: user-readable error message
    }

    return c.JSON(http.StatusCreated, CreateMessage("contract type successfully created")) //TODO: user-readable error message
}

func HandleGetAllContractType(c echo.Context) error {
    contractTypes, err := model.GetAllContractType()
    if err != nil {
        slog.Error("could not get all contract types", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrormessage("could not get all contract types")) //TODO: user-readable error message
    }
    return c.JSON(http.StatusOK, contractTypes)
}

func HandleUpdateContractType(c echo.Context) error {
    contractType := model.ContractType{}
    if err := c.Bind(&contractType); err != nil {
        return c.JSON(http.StatusInternalServerError, CreateErrormessage("could not bind contract type")) //TODO: user-readalbe error message
    }
    //TODO: decide if validating will be a thing

    if err := contractType.UpdateContractType(); err != nil {
        slog.Error("could not update contract type", "error", err, "contract type", contractType)
        return c.JSON(http.StatusBadRequest, CreateErrormessage("could not update contract type (possibly existing)")) //TODO: user-readable error message
    }

    return c.JSON(http.StatusOK, CreateMessage("contract type successfully updated")) //TODO: user-readable error message
}

func HandleDeleteContractType(c echo.Context) error {
    contractType := model.ContractType{}
    if err := c.Bind(&contractType); err != nil {
        return c.JSON(http.StatusInternalServerError, CreateErrormessage("could not bind contract type")) //TODO: user-readalbe error message
    }

    if err := contractType.DeleteContractType(); err != nil {
        slog.Error("could not delete contract type", "error", err, "contract type", contractType)
        return c.JSON(http.StatusBadRequest, CreateErrormessage("could not delete contract type")) //TODO: user-readable error message
    }

    return c.NoContent(http.StatusNoContent)
}

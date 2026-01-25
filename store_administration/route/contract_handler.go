package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

const (
    CONTRACT_FOLDER = "contracts"
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
    contractData := c.FormValue("contract")
	contract := model.Contract{}
	if err := json.Unmarshal([]byte(contractData), &contract); err != nil {
		slog.Error("could not bind contract", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés adatait nem sikerült értelmezni!"))
	}
    contract.Filename = sql.NullString{Valid: false}

	if err := contract.ValidateInsert(); err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
	}

    file, err := c.FormFile("file")
    if err == nil {
        if saveErr := saveContractFile(file); saveErr != nil {
            slog.Error("could not save contract file", "error", saveErr)
            return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not save contract file")) //TODO: hungarian error message
        }
        contract.Filename = sql.NullString{Valid: true, String: file.Filename}
    }

	if err := contract.InsertNewContract(); err != nil {
		slog.Error("could not insert new contract", "error", err, "contract", contract)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés mentése nem sikerült!"))
	}

	return c.JSON(http.StatusCreated, CreateMessage("A szerződés mentése sikeres!"))
}

func saveContractFile(file *multipart.FileHeader) error {
    src, err := file.Open()
    if err != nil {
        return err
    }

    dst, err := os.Create(fmt.Sprintf("%s/%s", CONTRACT_FOLDER, file.Filename))
    if err != nil {
        return err
    }

    _, err = io.Copy(dst, src)
    if err != nil {
        return err
    }

    return nil
}

func deleteContractFile(contract model.Contract) error {
    if !contract.Filename.Valid {
        return nil
    }
    return os.Remove(fmt.Sprintf("%s/%s", CONTRACT_FOLDER, contract.Filename.String))
}

func HandleUpdateContract(c echo.Context) error {
    contractData := c.FormValue("contract")
	contract := model.Contract{}
	if err := json.Unmarshal([]byte(contractData), &contract); err != nil {
		slog.Error("could not bind contract", "error", err)
		return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A szerződés adatait nem sikerült értelmezni!"))
	}

	if err := contract.ValidateUpdate(); err != nil {
		return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
	}

    file, err := c.FormFile("file")
    if err == nil {
        if delErr := deleteContractFile(contract); delErr != nil {
            slog.Error("could not delete previous contract file", "error", delErr)
            return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not delete previous contract file")) //TODO: hungarian error message
        }
        if saveErr := saveContractFile(file); saveErr != nil {
            slog.Error("could not save new contract file", "error", saveErr)
            return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not save new contract file")) //TODO: hungarian error message
        }
        contract.Filename = sql.NullString{Valid: true, String: file.Filename}
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

func HandleGetContractFile(c echo.Context) error {
    filename := c.QueryParam("filename")
    if filename == "" {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("missing filename")) //TODO: hungarian error message
    }
    return c.File(fmt.Sprintf("%s/%s", CONTRACT_FOLDER, filename))
}

func HandleDeleteContractFile(c echo.Context) error {
    contract := model.Contract{}
    if err := c.Bind(&contract); err != nil {
        slog.Error("could not bind contract for deleting contract file", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not bind contract")) //TODO: hungarian error message
    }
    if err := contract.DeleteContractFileFromDB(); err != nil {
        slog.Error("could not delete contract file from db", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not delete contract file from db")) //TODO: hungarian error message
    }
    if err := deleteContractFile(contract); err != nil {
        slog.Error("could not delete contract file", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not delete contract file")) //TODO: hungarian error message
    }
    return c.NoContent(http.StatusNoContent)
}

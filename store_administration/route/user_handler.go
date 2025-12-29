package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/auth"
	"github.com/vukovlevi/netstore/store_administration/model"
)

type PasswordUpdate struct {
    OldPassword string `json:"oldPassword"`
    NewPassword string `json:"newPassword"`
}

func HandlePostUser(c echo.Context) error {
    user := model.User{}
    if err := c.Bind(&user); err != nil {
        slog.Error("could not bind user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A felhasználó adatainak olvasása sikertelen!"))
    }

    if err := user.ValidateInsert(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    handlerUser := c.Get("user").(model.User)
    if !auth.CanUserSetRole(handlerUser, user.RoleId) {
        slog.Warn("a user tried to set a role that is higher then their own", "user who tried it", handlerUser)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("A jelenlegi rangoddal nincs jogod a kívánt rangot beállítani az új felhasználónak!"))
    }

    if user.PasswordChanged {
        if !auth.CanUserDisablePasswordChange(handlerUser) {
            slog.Warn("a user tried to disable password change, but doesnt have permission to do it", "user who tried it", handlerUser)
            return c.JSON(http.StatusBadRequest, CreateErrorMessage("A jelenlegi rangoddal nem kapcsolhatod ki a kötelező jelszó módosítást!"))
        }
    }

    if err := user.EncryptPassword(); err != nil {
        slog.Error("could not encrypt password of user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A felhasználó jelszavának titkosítása sikertelen, ezért nem kerül mentésre!"))
    }

    if err := user.InsertNewUser(); err != nil {
        slog.Error("could not save new user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A felhasználó mentése sikertelen (lehet, hogy már létezik a felhasználónév, akár a törölt felhasználók között is)!"))
    }

    return c.JSON(http.StatusCreated, CreateMessage("Az új felhasználó létrehozása sikeres!"))
}

func HandleUpdateUser(c echo.Context) error {
    user := model.User{}
    if err := c.Bind(&user); err != nil {
        slog.Error("could not bind user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A felhasználó adatainak olvasása sikertelen!"))
    }

    if err := user.ValidateUpdate(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    handlerUser := c.Get("user").(model.User)
    if !auth.CanUserSetRole(handlerUser, user.RoleId) {
        slog.Warn("a user tried to set a role that is higher then their own", "user who tried it", handlerUser)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("A jelenlegi rangoddal nem állíthatod be a kívánt rangra a felhasználót!"))
    }

    if err := user.UpdateUser(); err != nil {
        slog.Error("could not update user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A felhasználó módosítása sikertelen!"))
    }

    return c.JSON(http.StatusOK, CreateMessage("A felhasználó módosítása sikeres!"))
}

func HandleDeleteUser(c echo.Context) error {
    user := model.User{}
    if err := c.Bind(&user); err != nil {
        slog.Error("could not bind user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A felhasználó adatainak olvasása sikertelen!"))
    }

    deleteBy := c.Get("user").(model.User)

    if err := user.ValidateDelete(deleteBy, auth.ROLE_STORE_LEADER); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    if err := user.DeleteUser(); err != nil {
        slog.Error("could not delete user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A felhasználó törlése sikertelen!"))
    }

    return c.NoContent(http.StatusNoContent)
}

func HandleGetUser(c echo.Context) error {
    user := c.Get("user").(model.User)
    return c.JSON(http.StatusOK, user)
}

func HandleGetAllUser(c echo.Context) error {
    users, err := model.GetAllUser()
    if err != nil {
        slog.Error("could not get all user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A felhasználók lekérdezése sikertelen!"))
    }

    return c.JSON(http.StatusOK, users)
}

func HandleUpdateUserPassword(c echo.Context) error {
    newPassword := PasswordUpdate{}
    if err := c.Bind(&newPassword); err != nil {
        slog.Error("could not bind password update", "error", err)
        return c.JSON(http.StatusInternalServerError, "A jelszó módosításához szükséges adatok olvasása sikertelen!")
    }

    user := c.Get("user").(model.User)
    if err := user.ValidateUpdatePassword(newPassword.OldPassword, newPassword.NewPassword); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    if err := user.EncryptPassword(); err != nil {
        slog.Error("could not encrypt user password while updating said password", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A jelszó módosítása közben nem sikerült titkosítani azt, ezért nem került változtatásra!"))
    }

    if err := user.UpdatePassword(); err != nil {
        slog.Error("could not update password for user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A jelszó módosítása sikertelen!"))
    }

    return c.JSON(http.StatusOK, CreateMessage("A jelszó módosítása sikeres!"))
}

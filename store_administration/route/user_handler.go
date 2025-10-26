package route

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/auth"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func HandlePostUser(c echo.Context) error {
    user := model.User{}
    if err := c.Bind(&user); err != nil {
        slog.Error("could not bind user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not bind user")) //TODO: user-readable error message
    }

    if err := user.ValidateInsert(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    handlerUser := c.Get("user").(model.User)
    if !auth.CanUserSetRole(handlerUser, user.RoleId) {
        slog.Warn("a user tried to set a role that is higher then their own", "user who tried it", handlerUser)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("you are not allowed to set that role")) //TODO: user-readable error message
    }

    if user.PasswordChanged {
        if !auth.CanUserDisablePasswordChange(handlerUser) {
            slog.Warn("a user tried to disable password change, but doesnt have permission to do it", "user who tried it", handlerUser)
            return c.JSON(http.StatusBadRequest, CreateErrorMessage("you are not allowed to disable password change")) //TODO: user-readable error message
        }
    }

    if err := user.EncryptPassword(); err != nil {
        slog.Error("could not encrypt password of user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not encrypt user password")) //TODO: user-readable error message
    }

    if err := user.InsertNewUser(); err != nil {
        slog.Error("could not save new user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not save new user")) //TODO: user-readable error message
    }

    return c.JSON(http.StatusCreated, CreateMessage("new user successfully created"))
}

func HandleUpdateUser(c echo.Context) error {
    user := model.User{}
    if err := c.Bind(&user); err != nil {
        slog.Error("could not bind user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not bind user")) //TODO: user-readable error message
    }

    if err := user.ValidateUpdate(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    handlerUser := c.Get("user").(model.User)
    if !auth.CanUserSetRole(handlerUser, user.RoleId) {
        slog.Warn("a user tried to set a role that is higher then their own", "user who tried it", handlerUser)
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("you are not allowed to set that role")) //TODO: user-readable error message
    }

    if err := user.UpdateUser(); err != nil {
        slog.Error("could not update user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not update user")) //TODO: user-readable error message
    }

    return c.JSON(http.StatusOK, CreateMessage("user successfully updated")) //TODO: user-readable message
}

func HandleDeleteUser(c echo.Context) error {
    user := model.User{}
    if err := c.Bind(&user); err != nil {
        slog.Error("could not bind user", "error", err)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not bind user")) //TODO: user-readable error message
    }

    if err := user.ValidateDelete(); err != nil {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage(err.Error()))
    }

    if err := user.DeleteUser(); err != nil {
        slog.Error("could not delete user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not delete user")) //TODO: user-readable error message
    }

    return c.NoContent(http.StatusNoContent)
}

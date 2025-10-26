package route

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/auth"
	"github.com/vukovlevi/netstore/store_administration/model"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func HandleLogin(c echo.Context) error {
    loginReq := LoginRequest{}
    if err := c.Bind(&loginReq); err != nil {
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not bind login request")) //TODO: user-readable error message
    }

    if loginReq.Username == "" || loginReq.Password == "" {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("missing username or password")) //TODO: user-readable error message
    }

    ctx := context.Background()
    err := auth.LoginUser(loginReq.Username, loginReq.Password, &ctx)
    if err != nil {
        code, msg := createLoginErrorCodeAndMessage(err)
        return c.JSON(code, CreateErrorMessage(msg))
    }

    sessionContext := ctx.Value("session")
    session, ok := sessionContext.(model.Session)
    if !ok {
        slog.Error("could not read session from login context")
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not read session from login context")) //TODO: user-readable error message
    }

    SetSessionCookie(c, session)

    sessionUser := ctx.Value("user")
    user, ok := sessionUser.(model.User)
    if !ok {
        slog.Error("could not read user from login context")
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not read user from login context")) //TODO: user-readable error message
    }

    if !user.PasswordChanged {
        return c.Redirect(http.StatusSeeOther, "/password-change")
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "faca"})
}

func createLoginErrorCodeAndMessage(err error) (int, string) {
    switch err {
    case auth.ErrUserNotFound:
        return http.StatusBadRequest, "no user with username" //TODO: user-readable error message
    case auth.ErrBadPassword:
        return http.StatusBadRequest, "bad password" //TODO: user-readable error message
    default:
        slog.Error("could not log in user", "error", err)
        return http.StatusInternalServerError, "something went wrong during login" //TODO: user-readable error message
    }
}

func SetSessionCookie(c echo.Context, session model.Session) {
    c.SetCookie(&http.Cookie{
        Name: "auth_token",
        Value: session.Token,
        Path: "/",
        SameSite: http.SameSiteStrictMode,
        HttpOnly: true,
        Secure: false, //TODO: read from config
        Expires: session.ExpiresAt,
    })
}

package route

import (
	"context"
	"log/slog"
	"net/http"
	"time"

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

func HandleLogout(c echo.Context) error {
    user := c.Get("user").(model.User)
    if err := user.LogoutUser(); err != nil {
        slog.Error("could not log out user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("could not log out user")) //TODO: user-readable error message
    }


    c.SetCookie(&http.Cookie{
        Name: "auth_token",
        Value: "",
        Path: "/",
        SameSite: http.SameSiteStrictMode,
        HttpOnly: true,
        Secure: false, //TODO: read from config
        Expires: time.Now().Add(time.Second * -1),
        MaxAge: -1,
    })

    return c.JSON(http.StatusOK, CreateMessage("logout successful")) //TODO: user-readable message
}

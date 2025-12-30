package route

import (
	"context"
	"log/slog"
	"net/http"
	"os"
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
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A bejelentkezési adatok olvasása sikertelen!"))
    }

    if loginReq.Username == "" || loginReq.Password == "" {
        return c.JSON(http.StatusBadRequest, CreateErrorMessage("Hiányzó felhasználónév vagy jelszó!"))
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
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("A bejelentkezés során hiba lépett fel!"))
    }

    SetSessionCookie(c, session)
    return c.JSON(http.StatusOK, map[string]string{"message": "faca"}) //TODO: ezzel kezdeni valamit
}

func createLoginErrorCodeAndMessage(err error) (int, string) {
    switch err {
    case auth.ErrUserNotFound:
        return http.StatusBadRequest, "Hibás felhasználónév!"
    case auth.ErrBadPassword:
        return http.StatusBadRequest, "Hibás jelszó!"
    default:
        slog.Error("could not log in user", "error", err)
        return http.StatusInternalServerError, "A bejelentkezés során hiba lépett fel!"
    }
}

func SetSessionCookie(c echo.Context, session model.Session) {
    c.SetCookie(&http.Cookie{
        Name: "auth_token",
        Value: session.Token,
        Path: "/",
        SameSite: http.SameSiteStrictMode,
        HttpOnly: true,
        Secure: os.Getenv("SECURE_COOKIE") == "TRUE",
        Expires: session.ExpiresAt,
    })
}

func HandleLogout(c echo.Context) error {
    user := c.Get("user").(model.User)
    if err := user.LogoutUser(); err != nil {
        slog.Error("could not log out user", "error", err, "user", user)
        return c.JSON(http.StatusInternalServerError, CreateErrorMessage("Kijelentkezés sikertelen!"))
    }


    c.SetCookie(&http.Cookie{
        Name: "auth_token",
        Value: "",
        Path: "/",
        SameSite: http.SameSiteStrictMode,
        HttpOnly: true,
        Secure: os.Getenv("SECURE_COOKIE") == "TRUE",
        Expires: time.Now().Add(time.Second * -1),
        MaxAge: -1,
    })

    return c.JSON(http.StatusOK, CreateMessage("Sikeres kijelentkezés!"))
}

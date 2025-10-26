package middleware

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/model"
)

func AuthenticateUser(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        cookie, err := c.Cookie("auth_token")
        if err != nil {
            return c.Redirect(http.StatusTemporaryRedirect, "/login")
        }

        session, err := model.GetSessionByToken(cookie.Value)
        if err != nil {
            return c.Redirect(http.StatusTemporaryRedirect, "/login")
        }

        err = session.UpdateExpiry()
        if err != nil {
            slog.Error("could not update session expiry during authentication through middleware")
        }
        cookie.Expires = session.ExpiresAt
        c.SetCookie(cookie)

        user, err := model.GetUserByUserId(session.UserId)
        if err != nil {
            slog.Error("could not find user with id in the session", "error", err, "token", session.Token)
            cookie.MaxAge = -1
            c.SetCookie(cookie)
            return c.Redirect(http.StatusTemporaryRedirect, "/login")
        }

        if !user.PasswordChanged {
            return c.Redirect(http.StatusTemporaryRedirect, "/password-change")
        }

        c.Set("user", user)
        return next(c)
    }
}

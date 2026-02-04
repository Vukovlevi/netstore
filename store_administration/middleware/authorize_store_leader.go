package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/auth"
	"github.com/vukovlevi/netstore/store_administration/model"
	"github.com/vukovlevi/netstore/store_administration/route"
)

func AuthorizeStoreLeader(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        user := c.Get("user").(model.User)
        if user.Role != auth.ROLE_STORE_LEADER {
            return c.JSON(http.StatusUnauthorized, route.CreateErrorMessage("Nincs joga a kért művelethez!"))
        }
        return next(c)
    }
}

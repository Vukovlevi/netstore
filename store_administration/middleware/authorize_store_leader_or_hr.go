package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vukovlevi/netstore/store_administration/auth"
	"github.com/vukovlevi/netstore/store_administration/model"
	"github.com/vukovlevi/netstore/store_administration/route"
)

func AuthorizeStoreLeaderOrHR(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        user := c.Get("user").(model.User)
        if user.Role != auth.ROLE_STORE_LEADER && user.Role != auth.ROLE_HR {
            return c.JSON(http.StatusUnauthorized, route.CreateErrorMessage("unauthorized access")) //TODO: user-readable error message
        }
        return next(c)
    }
}

package forum

import (
	"github.com/labstack/echo/v4"

	cookie "forum/cookie"
)

func LogOut(c echo.Context) error {
        cookie.RemoveCookie(c)

        c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(200)
}

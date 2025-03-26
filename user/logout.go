package forum

import (
	"github.com/labstack/echo/v4"

	cookie "forum/cookie"
)

func LogOut(c echo.Context) error {
        cookie.RemoveCookie(c)

        c.Response().Header().Set("HX-Redirect", "/page/1")
	return c.NoContent(200)
}

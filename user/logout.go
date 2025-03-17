package forum

import (
        "github.com/labstack/echo/v4"

        cookie "forum/cookie"
)

func LogOut(c echo.Context) error {
        cookie.RemoveCookie(c)
        return c.Render(200, "home", nil)
}

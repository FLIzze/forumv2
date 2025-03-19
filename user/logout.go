package forum

import (
        "github.com/labstack/echo/v4"

        cookie "forum/cookie"
        structs "forum/structs"
)

func LogOut(c echo.Context) error {
        response := structs.User{}

        cookie.RemoveCookie(c)

        return c.Render(200, "navbar", response)
}

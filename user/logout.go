package forum

import (
        "github.com/labstack/echo/v4"

        cookie "forum/cookie"
        // structs "forum/structs"
)

func LogOut(c echo.Context) error {
        // response := structs.Status{}

        cookie.RemoveCookie(c)

        // response.Success = "Succesfully logout."

        return c.Render(200, "navbar", nil)
}

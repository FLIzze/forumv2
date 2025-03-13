package forum

import (
        "github.com/labstack/echo/v4"
)

func LogOut(c echo.Context) error {
        return c.Render(200, "login", nil)
}

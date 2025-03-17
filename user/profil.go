package forum

import (
        "github.com/labstack/echo/v4"
)

type Response struct {
        Error string
}

func Profil(c echo.Context) error {
        return c.Render(200, "profil", nil)
}

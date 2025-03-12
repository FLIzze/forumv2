package forum

import (
        "github.com/labstack/echo/v4"
)

func Get404(c echo.Context) error {
        return c.Render(404, "404", nil)
}

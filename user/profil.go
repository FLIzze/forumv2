package forum

import (
        "github.com/labstack/echo/v4"
        // "database/sql"
)

type Response struct {
        Error string
        User User
}

func Profil(c echo.Context) error {
        response := Response{} 

        // db := c.Get("db").(*sql.DB)
        // username := c.Param("username")

        return c.Render(200, "profil", response)
}

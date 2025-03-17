package forum

import (
        "github.com/labstack/echo/v4"
)

type Response struct {
        Error string
        User User
}

func Profil(c echo.Context) error {
        response := Response{} 

        user, ok := c.Get("user").(User)
        if !ok {
                c.Logger().Debug("User is not logged in")
        } else {
                response.User = user
        }

        return c.Render(200, "profil", response)
}

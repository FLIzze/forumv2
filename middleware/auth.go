package forum

import (
        "github.com/labstack/echo/v4"
        "database/sql"

        cookie "forum/cookie"
        user "forum/user"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                db, ok := c.Get("db").(*sql.DB)
                if !ok {
                        return echo.NewHTTPError(500, "Database connection not found")
                }

                cookie, err := cookie.GetCookie(c)
                if err == nil { 
                        var currentUser user.User
                        currentUser.SessionUUID = cookie.Value

                        row := db.QueryRow(`
                        SELECT 
                                UserUUID, Username, Email 
                        FROM 
                                userSession 
                        WHERE 
                                SessionUUID = ?
                        `, currentUser.SessionUUID)

                        err = row.Scan(&currentUser.UUID, &currentUser.Username, &currentUser.Email)
                        if err == nil {
                                c.Set("user", currentUser) 
                        }
                }

                if c.Path() == "/login" || c.Path() == "/register" {
                        return next(c)
                }

                if (c.Request().Method == echo.POST || c.Request().Method == echo.DELETE) && c.Get("user") == nil {
                        return echo.NewHTTPError(401, "You need to be logged in to perform this action.")
                }

                return next(c)
        }
}

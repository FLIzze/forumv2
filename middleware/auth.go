package forum

import (
        "github.com/labstack/echo/v4"
        "database/sql"

        cookie "forum/cookie"
        structs "forum/structs"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                db, ok := c.Get("db").(*sql.DB)
                if !ok {
                        c.Logger().Error("Error retrieving db from context.")
                        return echo.NewHTTPError(500, "Something went wrong. Please try again later.")
                }

                cookie, err := cookie.GetCookie(c)
                if err != nil { 
                        return next(c)
                }

                var currentUser structs.User
                currentUser.SessionUUID = cookie.Value

                row := db.QueryRow(`
                SELECT 
                        UserUUID
                FROM 
                        userSession 
                WHERE 
                        SessionUUID = ?
                `, currentUser.SessionUUID)

                err = row.Scan(&currentUser.UUID)
                if err != nil {
                        c.Logger().Error("Error retrieving UserUUID from userSession.", err)
                        return echo.NewHTTPError(500, "Something went wrong. Please try again later.")
                }

                row = db.QueryRow(`
                SELECT
                        Username, CreationTime, NmbMessagesPosted, NmbTopicsCreated, LastMessage
                FROM
                        userInfo
                WHERE 
                        UserUUID = ?
                `, currentUser.UUID)

                err = row.Scan(&currentUser.Username, &currentUser.CreationTime, &currentUser.NmbMessagesPosted, 
                                                        &currentUser.NmbTopicsCreated, &currentUser.LastMessage)
                if err != nil {
                        c.Logger().Error("Error retrieving userInfo.", err)
                        return echo.NewHTTPError(500, "Something went wrong. Please try again later.")
                }

                c.Set("user", currentUser) 

                if c.Path() == "/login" || c.Path() == "/register" {
                        return next(c)
                }

                if (c.Request().Method == echo.POST || c.Request().Method == echo.DELETE) && c.Get("user") == nil {
                        return c.HTML(401, `You must be logged in to perform this action. <a href="/">home</a> <a href="/login">login</a>`)
                }

                return next(c)
        }
}

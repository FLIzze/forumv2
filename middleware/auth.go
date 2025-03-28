package forum

import (
        "database/sql"

        "github.com/labstack/echo/v4"

        cookie "forum/cookie"
        structs "forum/structs"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                db, ok := c.Get("db").(*sql.DB)
                if !ok {
                        c.Logger().Error("Error retrieving db from context.")
                        return echo.NewHTTPError(500, "Something went wrong. Please try again later.")
                }

                cookie, err := cookie.GetCookie(c)
                if err != nil { 
                        c.Set("user", nil)
                        return next(c)
                }

                var currentUser structs.User
                currentUser.SessionUUID = cookie.Value

                row := db.QueryRow(`
                SELECT UserUUID
                FROM userSession 
                WHERE SessionUUID = ?
                `, currentUser.SessionUUID)

                err = row.Scan(&currentUser.UUID)
                if err != nil {
                        c.Set("user", nil)
                        return next(c)
                }

                row = db.QueryRow(`
                SELECT Username, CreationTime, NmbMessagesPosted, NmbTopicsCreated, LastMessage, LastTopic
                FROM userInfo
                WHERE UserUUID = ?
                `, currentUser.UUID)

                err = row.Scan(&currentUser.Username, &currentUser.CreationTime, &currentUser.NmbMessagesPosted, 
                                &currentUser.NmbTopicsCreated, &currentUser.LastMessage, &currentUser.LastTopic)
                if err != nil {
                        c.Logger().Error("Error retrieving userInfo.", err)
                        return echo.NewHTTPError(500, "Something went wrong. Please try again later.")
                }

                c.Set("user", currentUser) 

                return next(c)
        }
}

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                _, ok := c.Get("user").(structs.User)
                if !ok {
                        return c.Render(401, "unauthorized", nil)
                }

                return next(c)
        }
}

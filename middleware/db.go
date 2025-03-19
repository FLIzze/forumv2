package forum

import (
        "github.com/labstack/echo/v4"

        dbi "forum/db"
)

func DBMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
                db, err := dbi.ConnectDb()
                if err != nil {
                        c.Logger().Error("Error connecting to DB.")
                        return echo.NewHTTPError(500, "Something went wrong. Please try again later.")
                }
                defer db.Close()

                c.Set("db", db)

                err = dbi.CreateTable(db)
                if err != nil {
                        c.Logger().Error("Error creating table.")
                        return echo.NewHTTPError(500, "Something went wrong. Please try again later.")
                }

                err = dbi.CreateView(db)
                if err != nil {
                        c.Logger().Error("Error creating view.")
                        return echo.NewHTTPError(500, "Something went wrong. Please try again later.")
                }

                return next(c)
        }
}

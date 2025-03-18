package forum

import (
	"database/sql"

	"github.com/labstack/echo/v4"

        structs "forum/structs"
)

func GetDB(c echo.Context) (*sql.DB, error) {
        db, ok := c.Get("db").(*sql.DB)
        if !ok {
		return nil, echo.NewHTTPError(500, "Something went wrong. Please try again later.")
        }

        return db, nil
}

func GetUser(c echo.Context) (structs.User, error) {
	u, ok := c.Get("user").(structs.User)
	if !ok {
		return structs.User{}, echo.NewHTTPError(401, "You must be logged in")
	}

	return u, nil
}

package forum

import (
        "database/sql"
        "github.com/labstack/echo/v4"

        cookie "forum/cookie"
        utils "forum/utils"
        structs "forum/structs"
)

func GetLogin(c echo.Context) error {
        return c.Render(200, "login", nil)
}

func PostLogin(c echo.Context) error {
        response := structs.Status{}

        inputPassword := c.FormValue("password")
        inputUsername := c.FormValue("username")

        if (inputUsername == "" || inputPassword == "") {
                response.Error = "You must fill the whole form"
                return c.Render(422, "login-form", response)
        }

        db := c.Get("db").(*sql.DB)

        row := db.QueryRow(`
        SELECT Username, UUID, Password
        FROM user
        WHERE Username = ?
        `, inputUsername)

        var username string
        var uuid string
        var password []byte

        err := row.Scan(&username, &uuid, &password)
        if err != nil {
                c.Logger().Error("Incorrect username", err)
                response.Error = "Incorrect password or username"
                return c.Render(422, "login-form", response)
        }

        err = utils.CompareHashPassword(password, inputPassword)
        if err != nil {
                c.Logger().Error("Incorrect password", err)
                response.Error = "Incorrect password or username"
                return c.Render(422, "login-form", response)
        }

        sessionUUID := utils.Uuid()

        err = Login(db, uuid, sessionUUID, c)
        if err != nil {
                c.Logger().Error("Error updating session", err)
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "login-form", response)
        }

        c.Response().Header().Set("HX-Redirect", "/page/1")
	return c.NoContent(200)
}

func Login(db *sql.DB, uuid, sessionUUID string, c echo.Context) error {
        _, err := db.Exec(`
        UPDATE session
        SET Connected = 1, SessionUUID = ?
        WHERE UserUUID = ?
        `, sessionUUID, uuid)
        if err != nil {
                return err
        }

        cookie.PostCookie(c, sessionUUID)

        return nil
}

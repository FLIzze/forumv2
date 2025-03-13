package forum

import (
        "database/sql"
        "github.com/labstack/echo/v4"

        cookie "forum/cookie"
        utils "forum/utils"
)

type LoginResponse struct {
        Error string
        Login Login
}

type Login struct {
        Username string
        UUID string
        Password []byte
}

func GetLogin(c echo.Context) error {
        return c.Render(200, "login", nil)
}

func PostLogin(c echo.Context) error {
        response := LoginResponse{}

        password := c.FormValue("password")
        username := c.FormValue("username")

        db := c.Get("db").(*sql.DB)

        row := db.QueryRow(`
        SELECT Username, UUID, Password
        FROM user
        WHERE Username = ?
        `, username)

        err := row.Scan(&response.Login.Username, &response.Login.UUID, &response.Login.Password)
        if err != nil {
                c.Logger().Error("Incorrect username", err)
                response.Error = "Incorrect password or username"
                return c.Render(422, "login-form", response)
        }

        err = utils.CompareHashPassword(response.Login.Password, password)
        if err != nil {
                c.Logger().Error("Incorrect password", err)
                response.Error = "Incorrect password or username"
                return c.Render(422, "login-form", response)
        }

        sessionUUID := utils.Uuid()

        _, err = db.Exec(`
        UPDATE session
        SET Connected = 1, SessionUUID = ?
        WHERE UserUUID = ?
        `, sessionUUID, response.Login.UUID)
        if err != nil {
                c.Logger().Error("Error updating session", err)
                response.Error = "Internal server error"
                return c.Render(500, "login-form", response)
        }

        cookie.PostCookie(c, sessionUUID)

        return c.Render(200, "login-form", nil)
}

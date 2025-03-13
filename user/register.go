package forum

import (
	"database/sql"
	"github.com/labstack/echo/v4"

        utils "forum/utils"
)

type RegisterResponse struct {
        Error string
}

func GetRegister(c echo.Context) error {
        return c.Render(200, "register", nil)
}

func PostRegister(c echo.Context) error {
        response := RegisterResponse{}

        userUUID := utils.Uuid()
        password := c.FormValue("password")
        passwordConfirmation := c.FormValue("password-confirm")

        if (password != passwordConfirmation) {
                response.Error = "Both password must match"
                return c.Render(422, "register-form", response)
        }

        if (len(password) > 17 || len(password) < 3) {
                response.Error = "Password must be < 3 and > 17"
                return c.Render(422, "register-form", response)
        }

        hashedPassword, err := utils.GenerateHash(password)
        if err != nil {
                c.Logger().Error("Error hashing password: %s", err)
                response.Error = "Internal server error"
                return c.Render(500, "register-form", response)
        }

        var username = c.FormValue("username")

        if (len(username) < 3 || len(username) > 17) {
                response.Error = "Username must be < 3 and > 17"
                return c.Render(422, "register-form", response)
        }

        email := c.FormValue("email")

        db := c.Get("db").(*sql.DB)

        _, err = db.Exec(`
        INSERT INTO user (UUID, Username, Email, Password)
        VALUES (?, ?, ?, ?)
        `, userUUID, username, email, hashedPassword)
        if err != nil {
                c.Logger().Error("Error inserting user: %s", err)
                response.Error = "Internal server error"
                return c.Render(500, "register", response)
        }

        _, err = db.Exec(`
        INSERT INTO session (UserUUID)
        VALUES (?)
        `, userUUID)
        if err != nil {
                c.Logger().Error("Error inserting session: %s", err)
                response.Error = "Internal server error"
                return c.Render(500, "register", response)
        }

        return c.Render(200, "register", response)
}

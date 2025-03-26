package forum

import (
        "database/sql"
        "github.com/labstack/echo/v4"
        "time"

        utils "forum/utils"
        structs "forum/structs"
)

func GetRegister(c echo.Context) error {
        return c.Render(200, "register", nil)
}

func PostRegister(c echo.Context) error {
        response := structs.Status{}

        userUUID := utils.Uuid()
        username := c.FormValue("username")
        email := c.FormValue("email")
        password := c.FormValue("password")
        passwordConfirmation := c.FormValue("password-confirm")

        if (username == "" || email == "" || password == "" || passwordConfirmation == "") {
                response.Error = "You must fill the whole form"
                return c.Render(422, "register-form", response)
        }

        db := c.Get("db").(*sql.DB)

        rows, err := db.Query(`
        SELECT 
                Username, Email 
        FROM 
                user
        `)
        if err != nil {
                c.Logger().Error("Error retrieving username", err)
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "register-form", response)
        }
        defer rows.Close()

        for rows.Next() {
                var existingUsername string
                var existingEmail string

                err := rows.Scan(&existingUsername, &existingEmail)
                if err != nil {
                        c.Logger().Error("Error scanning row", err)
                        response.Error = "Something went wrong. Please try again later."
                        return c.Render(422, "register-form", response)
                }

                if username == existingUsername {
                        response.Error = "Username already taken"
                        return c.Render(422, "register-form", response)
                } else if email == existingEmail {
                        response.Error = "Email already taken"
                        return c.Render(422, "register-form", response)
                }
        }

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
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "register-form", response)
        }

        if (len(username) < 3 || len(username) > 17) {
                response.Error = "Username must be < 3 and > 17"
                return c.Render(422, "register-form", response)
        }

        _, err = db.Exec(`
        INSERT INTO user (UUID, Username, Email, Password, CreationTime)
        VALUES (?, ?, ?, ?, ?)
        `, userUUID, username, email, hashedPassword, time.Now())
        if err != nil {
                c.Logger().Error("Error inserting user: %s", err)
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "register", response)
        }

        sessionUUID := utils.Uuid()

        _, err = db.Exec(`
        INSERT INTO session (SessionUUID, UserUUID)
        VALUES (?, ?)
        `, sessionUUID, userUUID)
        if err != nil {
                c.Logger().Error("Error inserting session: %s", err)
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "register-form", response)
        }

        err = Login(db, userUUID, sessionUUID, c)
        if err != nil {
                c.Logger().Error("Error updating session", err)
                response.Error = "Something went wrong. Please try again later."
                return c.Render(500, "register-form", response)
        }

        c.Response().Header().Set("HX-Redirect", "/page/1")
	return c.NoContent(200)
}

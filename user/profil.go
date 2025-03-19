package forum

import (
        "database/sql"

        "github.com/labstack/echo/v4"

        structs "forum/structs"
)

func GetProfil(c echo.Context) error {
        response := structs.ProfilResponse{} 
        userProfil := structs.User{}

        user, ok := c.Get("user").(structs.User)
        if !ok {
                return c.HTML(401, `You must be logged in order to view a profil. <a href="/">home</a> <a href="/login">login</a>`)
        } else {
                response.User = user
        }

        db := c.Get("db").(*sql.DB)
        username := c.Param("username")

        row := db.QueryRow(`
        SELECT 
                Username, CreationTime, NmbMessagesPosted, NmbTopicsCreated, LastMessage
        FROM 
                userInfo
        WHERE
                Username = ?
        `, username)

        err := row.Scan(&userProfil.Username, &userProfil.CreationTime, &userProfil.NmbMessagesPosted, 
                                                &userProfil.NmbTopicsCreated, &userProfil.LastMessage)
        if err != nil {
                c.Logger().Error("Error retrieving user from userInfo (does not exist)", err)
                return c.Render(404, "404", nil)
        }

        response.UserProfil = userProfil
        return c.Render(200, "profil", response)
}

func GetMeProfil(c echo.Context) error {
        response := structs.ProfilResponse{}

        user, ok := c.Get("user").(structs.User)
        if !ok {
                return c.HTML(401, `You must be logged in order to view your profil. <a href="/">home</a> <a href="/login">login</a>`)
        } else {
                response.User = user
        }

        return c.Render(200, "meProfil", response)
}

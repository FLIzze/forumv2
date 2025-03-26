package forum

import (
        "database/sql"

        "github.com/labstack/echo/v4"

        structs "forum/structs"
        utils "forum/utils"
)

func GetProfil(c echo.Context) error {
        response := structs.ProfilResponse{} 
        userProfil := structs.User{}
        response.User = c.Get("user").(structs.User)

        db := c.Get("db").(*sql.DB)
        username := c.Param("username")

        row := db.QueryRow(`
        SELECT 
                Username, CreationTime, NmbMessagesPosted, NmbTopicsCreated, LastMessage, LastTopic
        FROM 
                userInfo
        WHERE
                Username = ?
        `, username)

        err := row.Scan(&userProfil.Username, &userProfil.CreationTime, &userProfil.NmbMessagesPosted, 
                        &userProfil.NmbTopicsCreated, &userProfil.LastMessage, &userProfil.LastTopic)
        if err != nil {
                c.Logger().Error("Error retrieving user from userInfo (does not exist)", err)
                return c.Render(404, "404", nil)
        }

        userProfil.FormattedCreationTime = utils.FormatDate(userProfil.CreationTime)
        response.UserProfil = userProfil
        return c.Render(200, "profil", response)
}

func GetMeProfil(c echo.Context) error {
        response := structs.ProfilResponse{}
        response.User = c.Get("user").(structs.User)

        return c.Render(200, "meProfil", response)
}

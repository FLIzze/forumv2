package forum

import (
        "database/sql"

        "github.com/labstack/echo/v4"

        structs "forum/structs"
)

func GetProfil(c echo.Context) error {
        response := structs.ProfilResponse{} 
        user := structs.User{}

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(structs.User)
        if !ok {
                response.Status.Error = "You must be logged in to delete a topic."
                return c.Render(401, "home", response)
        } else {
                response.User = user
        }
        username := c.Param("username")

        row := db.QueryRow(`
        SELECT 
                Username, CreationTime, NmbMessagesPosted, NmbTopicsCreated, LastMessage
        FROM 
                userInfo
        WHERE
                Username = ?
        `, username)

        err := row.Scan(&user.Username, &user.CreationTime, &user.NmbMessagesPosted, &user.NmbTopicsCreated, &user.LastMessage)
        if err != nil {
                c.Logger().Error("Error retrieving user from userInfo", err)
                response.Status.Error = "User does not exist"
                return c.Render(404, "404", nil)
        }
        response.UserProfil = user

        return c.Render(200, "profil", response)
}

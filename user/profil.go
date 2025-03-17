package forum

import (
        "github.com/labstack/echo/v4"
        "database/sql"
)

type Response struct {
        Error string
        User User
}

func Profil(c echo.Context) error {
        response := Response{} 
        user := User{}

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

        err := row.Scan(&user.Username, &user.CreationTime, &user.NmbMessagesPosted, &user.NmbTopicsCreated, &user.LastMessage)
        if err != nil {
                c.Logger().Error("Error retrieving user from userInfo", err)
                response.Error = "User does not exist"
                return c.Render(404, "404", nil)
        }
        response.User = user

        return c.Render(200, "profil", response)
}

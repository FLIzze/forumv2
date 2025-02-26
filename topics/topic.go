package forum

import (
	"github.com/labstack/echo/v4"
        dbi "forum/db"
)

func GetTopic(c echo.Context) error {
        topic := Topic{}
        URI := c.Request().RequestURI[1:]
        db := dbi.ConnectDb()

        rows, err := db.Query(`
                SELECT UUID, Name, Description 
                FROM topic
                WHERE UUID = ?
        `, URI)
        if err != nil {
                panic(err)
        }
        defer rows.Close()

        for rows.Next() {
                var uuid, name, description string

                err := rows.Scan(&uuid, &name, &description)
                if err != nil {
                        panic(err)
                }

                topic = NewTopic(uuid, name, description)
        }

        return c.Render(200, "topic.html", topic)
}

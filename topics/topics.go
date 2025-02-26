package forum

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
        dbi "forum/db"
)

func PostTopic(c echo.Context) error {
        db := dbi.ConnectDb()
        defer db.Close()

        uuid := uuid.New().String()
        name := c.FormValue("name")
        description := c.FormValue("description")

        _, err := db.Exec(`
                INSERT INTO topic (UUID, Name, Description) 
                VALUES (?, ?, ?)
        `, uuid, name, description)
        if err != nil {
                panic(err)
        }

        newTopic := NewTopic(uuid, name, description)

        return c.Render(200, "oob-topic", newTopic)
}

func GetTopics(c echo.Context) error {
        db := dbi.ConnectDb()
        defer db.Close()

        topics := Topics{}

        rows, err := db.Query(`
                SELECT UUID, Name, Description FROM topic
        `)
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
                
                topics.Topics = append(topics.Topics, NewTopic(uuid, name, description))
        }

        return c.Render(200, "index", topics)
}

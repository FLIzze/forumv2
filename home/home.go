package forum

import (
        "database/sql"
        "github.com/google/uuid"
        "github.com/labstack/echo/v4"

        dbi "forum/db"
)

type HomeResponse struct {
        Error string
        Topics []Topic
}

type Topic struct {
        UUID string
        Name string
        Description string
}

func HomePage(c echo.Context) error {
        response := HomeResponse{}
        topic := Topic{}

        db := c.Get("db").(*sql.DB)

        rows, err := db.Query(`
        SELECT UUID, Name, Description FROM topic
        `)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Error = "Could not retrieve topic."
                return c.Render(500, "home", response)
        }
        defer rows.Close()

        for rows.Next() {
                err := rows.Scan(&topic.UUID, &topic.Name, &topic.Description)
                if err != nil {
                        c.Logger().Error("Error scanning row", err)
                        response.Error = "Could not retrieve data from column"
                        return c.Render(422, "home", response)
                }

                response.Topics = append(response.Topics, topic)
        }

        return c.Render(200, "home", response)
}

func PostTopic(c echo.Context) error {
        response := HomeResponse{}
        topic := Topic{}

        topic.UUID = uuid.New().String()
        topic.Name = c.FormValue("name")
        topic.Description = c.FormValue("description")

        if topic.Name == "" || topic.Description == "" {
                c.Logger().Error("Name and/or description empty")
                response.Error = "Name and/or description must be filled."
                return c.Render(422, "topics-form", response)
        }

        response.Topics = append(response.Topics, topic)

        db := c.Get("db").(*sql.DB)

        _, err := db.Exec(`
        INSERT INTO topic (UUID, Name, Description) 
        VALUES (?, ?, ?)
        `, topic.UUID, topic.Name, topic.Description)
        if err != nil {
                c.Logger().Error("Error retrieving response.Topic: ", err)
                response.Error = "Error retrieving response.Topic."
                return c.Render(500, "topics-form", response)
        }

        c.Render(200, "topics-form", response)
        return c.Render(200, "oob-topic", response)
}

func DeleteTopic(c echo.Context) error {
        URI := c.Param("uuid")

        db, err := dbi.ConnectDb()
        if err != nil {
                c.Logger().Error("Error connecting to db: ", err)
                return c.Render(500, "error", nil)
        }

        _, err = db.Exec(`
        DELETE FROM topic WHERE UUID = ?
        `, URI)
        if err != nil {
                c.Logger().Error("Error deleting topic: ", err)
                return c.Render(500, "error", nil)
        }

        return c.Render(200, "oob-topic", nil)
}

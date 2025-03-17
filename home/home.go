package forum

import (
        "database/sql"

        "github.com/google/uuid"
        "github.com/labstack/echo/v4"

        user "forum/user"
)

type HomeResponse struct {
        Error string
        Topics []Topic
        User user.User
}

type Topic struct {
        UUID string
        Name string
        Description string
        Index int
}

func GetHomePage(c echo.Context) error {
        response := HomeResponse{}
        topic := Topic{}

        db := c.Get("db").(*sql.DB)
        user, ok := c.Get("user").(user.User)
        if ok {
                response.User = user
        }

        rows, err := db.Query(`
        SELECT UUID, Name, Description FROM topic
        `)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Error = "Could not retrieve topic."
                return c.Render(500, "home", response)
        }
        defer rows.Close()

        index := 0

        for rows.Next() {
                err := rows.Scan(&topic.UUID, &topic.Name, &topic.Description)
                if err != nil {
                        c.Logger().Error("Error scanning row", err)
                        response.Error = "Internal error"
                        return c.Render(422, "home", response)
                }

                topic.Index = index
                response.Topics = append(response.Topics, topic)
                index++
        }

        c.Logger().Error(user)

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
        user, ok := c.Get("user").(user.User)
        if !ok {
                c.Logger().Debug("User is not logged in")
        } else {
                response.User = user
        }

        _, err := db.Exec(`
        INSERT INTO topic (UUID, Name, Description, CreatedBy) 
        VALUES (?, ?, ?, ?)
        `, topic.UUID, topic.Name, topic.Description, user.UUID)
        if err != nil {
                c.Logger().Error("Error inserting response Topic: ", err)
                response.Error = "Internal error"
                return c.Render(500, "topics-form", response)
        }

        c.Render(200, "topics-form", response)
        return c.Render(200, "oob-topic", response)
}

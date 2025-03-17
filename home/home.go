package forum

import (
	"database/sql"
	"time"

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
        CreatedBy string
        NmbMessages int
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
        SELECT 
                UUID, Name, Description, CreatedBy, NmbMessages 
        FROM 
                topicInfo
        `)
        if err != nil {
                c.Logger().Error("Error retrieving topic: ", err)
                response.Error = "Could not retrieve topic."
                return c.Render(500, "home", response)
        }
        defer rows.Close()

        for rows.Next() {
                err := rows.Scan(&topic.UUID, &topic.Name, &topic.Description, &topic.CreatedBy, &topic.NmbMessages)
                if err != nil {
                        c.Logger().Error("Error scanning row", err)
                        response.Error = "Internal server error"
                        return c.Render(422, "home", response)
                }

                response.Topics = append(response.Topics, topic)
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
        INSERT INTO topic (UUID, Name, Description, CreatedBy, CreationTime) 
        VALUES (?, ?, ?, ?, ?)
        `, topic.UUID, topic.Name, topic.Description, user.UUID, time.Now())
        if err != nil {
                c.Logger().Error("Error inserting response Topic: ", err)
                response.Error = "Internal error"
                return c.Render(500, "topics-form", response)
        }

        c.Render(200, "topics-form", response)
        return c.Render(200, "oob-topic", response)
}
